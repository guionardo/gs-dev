package cmdtools

import (
	"github.com/guionardo/gs-dev/internal"
	"github.com/spf13/cobra"
)

type CmdParser struct {
	args       []string
	stringArgs internal.Set
	intArgs    internal.Set
	boolArgs   internal.Set
}

func CreateCmdParser() *CmdParser {
	return &CmdParser{
		args:       []string{},
		stringArgs: internal.NewSet(),
		intArgs:    internal.NewSet(),
		boolArgs:   internal.NewSet(),
	}
}
func (parser *CmdParser) StringArg(argName string) *CmdParser {
	parser.stringArgs.Add(argName)
	return parser
}
func (parser *CmdParser) IntArg(argName string) *CmdParser {
	parser.intArgs.Add(argName)
	return parser
}

func (parser *CmdParser) GetCmdArgs(cmd *cobra.Command) map[string]interface{} {
	args := make(map[string]interface{})

	for _, key := range parser.stringArgs.ToList() {
		if value, err := cmd.Flags().GetString(key); err == nil {
			args[key] = value
		}
	}
	for _, key := range parser.intArgs.ToList() {
		if value, err := cmd.Flags().GetInt(key); err == nil {
			args[key] = value
		}
	}
	for _, key := range parser.boolArgs.ToList() {
		if value, err := cmd.Flags().GetBool(key); err == nil {
			args[key] = value
		}
	}
	return args
}
