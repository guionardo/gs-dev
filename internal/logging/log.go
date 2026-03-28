package logging

import (
	"fmt"
	"io"
	"log/slog"
)

type preSetupLog struct {
	level slog.Level
	msg   string
}

var preSetupLogs []preSetupLog

func PreSetupLog(msg string, level slog.Level, args ...any) {
	preSetupLogs = append(preSetupLogs, preSetupLog{level: level, msg: fmt.Sprintf(msg, args...)})
}

func Setup(debug bool, output io.Writer) {
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if debug {
		handlerOptions.Level = slog.LevelDebug
		handlerOptions.AddSource = true
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(output, handlerOptions)))

	for _, log := range preSetupLogs {
		switch log.level {
		case slog.LevelDebug:
			slog.Debug(log.msg)
		case slog.LevelWarn:
			slog.Warn(log.msg)
		case slog.LevelError:
			slog.Error(log.msg)
		default:
			slog.Info(log.msg)
		}
	}

	preSetupLogs = nil
}
