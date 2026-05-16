package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/guionardo/go/flow"
	timetools "github.com/guionardo/go/time_tools"
	"github.com/guionardo/gs-dev/internal/consts"
	errs "github.com/guionardo/gs-dev/internal/errors"
	string_tools "github.com/guionardo/gs-dev/internal/tools/strings_tools"
	"github.com/spf13/cobra"
)

type (
	// CommandStruct is the contract used by GenerateCobraCommand.
	//
	// Implementations are usually pointer-to-struct command definitions whose
	// exported fields are populated from Cobra flags/args before Run is called.
	CommandStruct interface {
		Run(ctx context.Context, output io.Writer) error
		Setup(ctx context.Context) error
	}

	flagMetadata struct {
		fieldIndex   int
		name         string
		shortHand    string
		defaultValue string
		description  string
		persistent   bool
		required     bool
		validate     bool
		stdin        bool
	}
	mapFlags map[int]flagMetadata
)

var layouts = []string{
	time.DateTime,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateOnly,
	time.TimeOnly,
}

// GenerateCobraCommand creates a cobra.Command from a command definition struct.
//
// Supported field tags in exported struct fields:
//   - `flag:"name[,short]"` to expose a flag.
//   - `description:"text"` to describe the flag.
//   - `default:"value"` to set a flag default value.
//   - `validate:"..."` to enable go-playground/validator checks before Run.
//   - `subcommand:"name,usage"` on struct fields that implement CommandStruct.
//   - `args:"..."` on []string fields to receive positional arguments.
//
// Supported flag field types: string, int, bool, []string and time.Duration.
// Any unsupported type or invalid default value causes a panic.
func GenerateCobraCommand(instance CommandStruct, name, usage, longDescription string, useOutput bool) *cobra.Command {
	t := validateTypeCobraStruct(instance)

	var (
		tag               string
		ok                bool
		flagsMap          = make(mapFlags) // flagsMap maps field index to flag name
		argsSliceIndex    = -1
		argsExpectedCount = 0
		argsDescription   = ""
	)

	cmd := &cobra.Command{
		Use:   name,
		Short: usage,
		Long:  longDescription,
	}

	for fieldIndex := range t.NumField() {
		field := t.Field(fieldIndex)
		if !field.IsExported() {
			continue
		}

		if parseSubcommand(cmd, field) {
			continue
		}

		var flag flagMetadata

		_, flag.validate = field.Tag.Lookup("validate")
		flag.required = strings.Contains(field.Tag.Get("validate"), "required")

		if tag, ok = field.Tag.Lookup("flag"); ok {
			words := strings.Split(tag, ",")
			flag.fieldIndex = fieldIndex

			flag.name = words[0]
			if len(words) > 1 {
				flag.shortHand = words[1]
			}

			flag.description = field.Tag.Get("description")
			_, flag.persistent = field.Tag.Lookup("persistent")
			flag.defaultValue = field.Tag.Get("default")
			addFlag(t.Name(), cmd, &field, flag)

			flagsMap[fieldIndex] = flag

			continue
		}

		// stdin
		found, err := parseStdinTag(field)
		if err != nil {
			panic(fmt.Sprintf("GenerateCobraCommand: %v", err))
		}

		if found {
			flagsMap[fieldIndex] = flagMetadata{
				fieldIndex: fieldIndex,
				stdin:      true,
			}

			continue
		}
		// Args
		if found, expectedCount, description, err := parseArgsTag(field); found {
			if err != nil {
				panic(fmt.Sprintf("GenerateCobraCommand: %v", err))
			}

			argsExpectedCount = expectedCount
			argsDescription = description
			argsSliceIndex = fieldIndex
		}
	}

	if argsExpectedCount == -1 {
		cmd.Args = cobra.MinimumNArgs(1)
	} else if argsExpectedCount >= 0 {
		cmd.Args = cobra.ExactArgs(argsExpectedCount)
	}

	if argsDescription != "" {
		cmd.Long += "\n\nArguments:\n" + argsDescription
	}

	cmd.RunE = createRunCommand(instance, flagsMap, argsSliceIndex)

	annotations := map[string]string{}
	if useOutput {
		annotations[consts.UseOutputAnnotation] = "true"
	}

	cmd.Annotations = annotations

	return cmd
}

func parseStdinTag(field reflect.StructField) (found bool, err error) {
	if stdin := field.Tag.Get("stdin"); stdin == "" {
		return false, nil
	}
	// Check if the field is a slice of bytes
	if field.Type.Kind() != reflect.Slice || field.Type.Elem().Kind() != reflect.Uint8 {
		err = fmt.Errorf("field %s.%s has an 'stdin' tag but is not a slice of bytes", field.Type.Name(), field.Name)
		return
	}

	return true, nil
}

// parseArgsTag parses the `args` struct tag and returns the expected argument count and description.
// The `args` tag format is either "X,description"
// X can be a number to expect exactly N arguments, -1 if at least 1 argument is expected.
// The description is required and is used in the command's Long description.
func parseArgsTag(field reflect.StructField) (found bool, expectedCount int, description string, err error) {
	args := field.Tag.Get("args")
	// Check if the args tag is present
	if args == "" {
		return
	}

	// Check if the field is a slice of strings
	if field.Type.Kind() != reflect.Slice || field.Type.Elem().Kind() != reflect.String {
		err = fmt.Errorf("field %s.%s has an 'args' tag but is not a slice of strings", field.Type.Name(), field.Name)
		return
	}

	// Parse the args tag

	var nArgs, argsDescription string
	string_tools.SplitString(args, ",", &nArgs, &argsDescription)

	if expectedCount, err = strconv.Atoi(nArgs); err == nil {
		if argsDescription == "" {
			err = errors.New("missing description")
		}
	}

	if err != nil {
		err = fmt.Errorf("invalid args tag format for field %s.%s: %w", field.Type.Name(), field.Name, err)
		return
	}

	description = argsDescription
	found = true

	return
}

func parseSubcommand(cmd *cobra.Command, field reflect.StructField) bool {
	// Check the field is not a struct, it can't be a subcommand
	if field.Type.Kind() != reflect.Struct {
		return false
	}
	// Check if the field has a tag `flag`
	if _, isFlag := field.Tag.Lookup("flag"); isFlag {
		return false
	}

	// Check if the field has a tag `subcommand:"name,usage"`
	subCommandTag := field.Tag.Get("subcommand")
	if subCommandTag == "" {
		return true // if it's a struct but doesn't have the subcommand tag, we consider it parsed to avoid further processing
	}

	var subCommandName, subCommandUsage, longDescription string
	string_tools.SplitString(subCommandTag, ",", &subCommandName, &subCommandUsage, &longDescription)

	// Check if the struct field implements CommandStruct
	// commandStructType := reflect.TypeFor[CommandStruct]()

	// if !field.Type.Implements(commandStructType) {
	// 	return true // if it doesn't implement CommandStruct, we consider it parsed to avoid further processing
	// }

	subCommandInstance := reflect.New(field.Type)

	sciInterface := subCommandInstance.Interface()
	if sci, ok := (sciInterface).(CommandStruct); ok {
		subCommand := GenerateCobraCommand(sci, subCommandName, subCommandUsage, longDescription, false)
		cmd.AddCommand(subCommand)
	}

	return true
}

func addFlag(typeName string, cmd *cobra.Command, field *reflect.StructField, flag flagMetadata) {
	flagSet := flow.If(flag.persistent, cmd.PersistentFlags(), cmd.Flags())

	if flag.required {
		defer func() {
			cmd.MarkFlagRequired(flag.name)
		}()
	}

	switch field.Type.Kind() {
	case reflect.String:
		if flag.shortHand == "" {
			flagSet.String(flag.name, flag.defaultValue, flag.description)
		} else {
			flagSet.StringP(flag.name, flag.shortHand, flag.defaultValue, flag.description)
		}

		return
	case reflect.Int:
		defaultValue, err := strconv.Atoi(flag.defaultValue)
		if flag.defaultValue != "" && err != nil {
			panic(fmt.Sprintf("GenerateCobraCommand.addFlag '%s' with default value '%s' is not a int value [%s]", flag.name, flag.defaultValue, typeName))
		}

		if flag.shortHand == "" {
			flagSet.Int(flag.name, defaultValue, flag.description)
		} else {
			flagSet.IntP(flag.name, flag.shortHand, defaultValue, flag.description)
		}

		return
	case reflect.Bool:
		defaultValue, err := strconv.ParseBool(flag.defaultValue)
		if flag.defaultValue != "" && err != nil {
			panic(fmt.Sprintf("GenerateCobraCommand.addFlag '%s' with default value '%s' is not a bool value [%s]", flag.name, flag.defaultValue, typeName))
		}

		if flag.shortHand == "" {
			flagSet.Bool(flag.name, defaultValue, flag.description)
		} else {
			flagSet.BoolP(flag.name, flag.shortHand, defaultValue, flag.description)
		}

		return
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.String {
			defaultValue := strings.Split(flag.defaultValue, ",")
			if flag.shortHand == "" {
				flagSet.StringSlice(flag.name, defaultValue, flag.description)
			} else {
				flagSet.StringSliceP(flag.name, flag.shortHand, defaultValue, flag.description)
			}

			return
		}
	case reflect.Int64:
		if field.Type.String() == "time.Duration" {
			defaultValue, err := time.ParseDuration(flag.defaultValue)
			if flag.defaultValue != "" && err != nil {
				panic(fmt.Sprintf("GenerateCobraCommand.addFlag '%s' with default value '%s' is not a time.Duration value [%s]", flag.name, flag.defaultValue, typeName))
			}

			if flag.shortHand == "" {
				flagSet.Duration(flag.name, defaultValue, flag.description)
			} else {
				flagSet.DurationP(flag.name, flag.shortHand, defaultValue, flag.description)
			}

			return
		}

		panic(fmt.Sprintf("unexpected flag field with type %s for command struct %s", field.Type.Name(), typeName))

	case reflect.Struct:
		if field.Type.String() == "time.Time" {
			defaultValue, err := timetools.Parse(flag.defaultValue)
			if flag.defaultValue != "" && err != nil {
				panic(fmt.Sprintf("GenerateCobraCommand.addFlag '%s' with default value '%s' is not a time.Time value [%s]", flag.name, flag.defaultValue, typeName))
			}

			if flag.shortHand == "" {
				flagSet.Time(flag.name, defaultValue, layouts, flag.description)
			} else {
				flagSet.TimeP(flag.name, flag.shortHand, defaultValue, layouts, flag.description)
			}

			return
		}
	default:
		panic(fmt.Sprintf("unexpected flag field with type %s for command struct %s", field.Type.Name(), typeName))
	}
}

func parseArgsAndFlags(instance CommandStruct, cmd *cobra.Command, args []string, flagsmap mapFlags, argsSliceIndex int) error {
	var (
		valueString string
		valueInt    int
		valueBool   bool
		err         error
	)

	csValue := reflect.ValueOf(instance).Elem()

	// Args
	if argsSliceIndex >= 0 {
		argsField := csValue.Field(argsSliceIndex)
		argsValue := reflect.ValueOf(args)
		argsField.Set(argsValue)
	}

	// Flags
	for index, flag := range flagsmap {
		flagField := csValue.Field(index)

		switch flagField.Type().Kind() {
		case reflect.String:
			if valueString, err = cmd.Flags().GetString(flag.name); err == nil {
				flagField.SetString(valueString)
			}

		case reflect.Int:
			if valueInt, err = cmd.Flags().GetInt(flag.name); err == nil {
				flagField.SetInt(int64(valueInt))
			}

		case reflect.Bool:
			if valueBool, err = cmd.Flags().GetBool(flag.name); err == nil {
				flagField.SetBool(valueBool)
			}

		case reflect.Int64:
			switch flagField.Type().String() {
			case "time.Duration":
				if valueDuration, err := cmd.Flags().GetDuration(flag.name); err == nil {
					valueDurationVal := reflect.ValueOf(valueDuration)
					flagField.Set(valueDurationVal)
				}
			}
		case reflect.Struct:
			switch flagField.Type().String() {
			case "time.Time":
				if valueTime, err := cmd.Flags().GetTime(flag.name); err == nil {
					valueTimeVal := reflect.ValueOf(valueTime)
					flagField.Set(valueTimeVal)
				}
			}

		case reflect.Slice:
			switch flagField.Type().Elem().Kind() {
			case reflect.String:
				if valueSlice, err := cmd.Flags().GetStringSlice(flag.name); err == nil {
					valueSliceVal := reflect.ValueOf(valueSlice)
					flagField.Set(valueSliceVal)
				}
			case reflect.Uint8:
				// reader := cmd.InOrStdin()
				stdInData, err := ReadFromStdIn()
				if err == nil {
					valueSliceVal := reflect.ValueOf(stdInData)
					flagField.Set(valueSliceVal)
				}
			}

		default:
			err = fmt.Errorf("unexpected field type for %T.%s", instance, flagField.Type().Name())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func createRunCommand(instance CommandStruct, flagsmap mapFlags, argsSliceIndex int) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := instance.Setup(cmd.Context()); err != nil {
			return err
		}

		err := parseArgsAndFlags(instance, cmd, args, flagsmap, argsSliceIndex)
		if err != nil {
			return err
		}

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		if flagsmap.hasValidation() {
			val := validator.New()
			if err = val.StructCtx(ctx, instance); err != nil {
				return err
			}
		}

		return instance.Run(ctx, cmd.OutOrStdout())
	}
}

func (mf mapFlags) hasValidation() bool {
	for _, flag := range mf {
		if flag.validate {
			return true
		}
	}

	return false
}

func validateTypeCobraStruct(instance CommandStruct) reflect.Type {
	t := reflect.TypeOf(instance)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		panic("GenerateCobraCommand expects a struct type. Got " + t.String() + " [" + t.Kind().String() + "]")
	}

	return t.Elem()
}

func ReadFromStdIn() (content []byte, err error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped
		content, err = io.ReadAll(os.Stdin)
	} else {
		err = errs.NewError(nil, "no data in stdin", false)
	}

	return
}
