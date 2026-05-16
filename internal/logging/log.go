package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type preSetupLog struct {
	level slog.Level
	msg   string
	args  []any
}

var (
	preSetupLogs []preSetupLog
	logDebug     bool
	logger       *slog.Logger
	mu           sync.Mutex
)

const moduleName = "github.com/guionardo/gs-dev/"

// Logger get/create the default logger. Use args to define level "--debug" and io.Writer to customize
func Logger(args ...any) *slog.Logger {
	mu.Lock()
	defer mu.Unlock()

	if logger == nil || len(args) > 0 {
		level, writer := logSetupArgs(args...)
		handlerOptions := &slog.HandlerOptions{
			Level: level,
		}
		logDebug = level == slog.LevelDebug

		logger = slog.New(slog.NewTextHandler(writer, handlerOptions))
		logger.Debug("Logging started")
	}

	if !logDebug {
		return logger
	}

	var source, callerFunc string

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		source = fmt.Sprintf("%s:%d", file, line)
		callerFunc = runtime.FuncForPC(pc).Name()

		return logger.With(slog.String("source", source), slog.String("caller", strings.TrimPrefix(callerFunc, moduleName)))
	}

	return logger
}

func Debug(msg string, args ...any) {
	Logger().Debug(msg, args...)
}

func Error(msg string, args ...any) {
	Logger().Error(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger().Warn(msg, args...)
}

func Info(msg string, args ...any) {
	Logger().Info(msg, args...)
}

func logSetupArgs(args ...any) (level slog.Level, writer io.Writer) {
	level = slog.LevelInfo
	writer = os.Stdout

	if len(args) == 0 { // Get args from command line
		for _, arg := range os.Args[1:] {
			args = append(args, arg)
		}
	}

	for _, arg := range args {
		if s, ok := arg.(string); ok {
			if s == "--debug" {
				level = slog.LevelDebug
			}
		}

		if b, ok := arg.(bool); ok {
			if b {
				level = slog.LevelDebug
			}
		}

		if w, ok := arg.(io.Writer); ok {
			writer = w
		}
	}

	return level, writer
}

// SetupDebug check the program args and looks for --debug flag
func SetupDebug(args ...string) {
	if len(args) == 0 {
		args = os.Args[1:]
	}

	logDebug = slices.Contains(args, "--debug")
}

func PreSetupLog(msg string, level slog.Level, args ...any) {
	// Get caller
	var source, callerFunc string

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		source = fmt.Sprintf("%s:%d", file, line)
		callerFunc = runtime.FuncForPC(pc).Name()

		if strings.Contains(msg, "%") {
			msg = fmt.Sprintf(msg, args...)
			args = []any{}
		}

		args = append(args, slog.String("source", source), slog.String("caller", callerFunc))
	}

	mu.Lock()

	preSetupLogs = append(preSetupLogs, preSetupLog{level: level, msg: msg, args: args})
	mu.Unlock()
}

func Setup(output io.Writer) {
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if logDebug {
		handlerOptions.Level = slog.LevelDebug
		// handlerOptions.AddSource = true
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(output, handlerOptions)))

	var logFn func(msg string, args ...any)

	for _, log := range preSetupLogs {
		switch log.level {
		case slog.LevelDebug:
			logFn = slog.Debug
			slog.Debug(log.msg)
		case slog.LevelWarn:
			logFn = slog.Warn
			slog.Warn(log.msg)
		case slog.LevelError:
			logFn = slog.Error
			slog.Error(log.msg)
		default:
			logFn = slog.Info
			slog.Info(log.msg)
		}

		if strings.Contains(log.msg, "%") {
			// Log using fmt
			logFn(fmt.Sprintf(log.msg, log.args...))
		} else {
			// Log using slog args
			logFn(log.msg, log.args...)
		}
	}

	preSetupLogs = nil
}

func getCaller(args ...any) []any {
	if !logDebug {
		return args
	}

	var source, callerFunc string

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		source = fmt.Sprintf("%s:%d", file, line)
		callerFunc = runtime.FuncForPC(pc).Name()

		args = append(args, slog.String("source", source), slog.String("caller", callerFunc))
	}

	return args
}
