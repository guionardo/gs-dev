package logging_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/guionardo/gs-dev/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Parallel()

	output := bytes.NewBufferString("")

	logging.PreSetupLog("Pre setup log - debug", slog.LevelDebug)
	logging.PreSetupLog("Pre setup log - info", slog.LevelInfo)
	logging.PreSetupLog("Pre setup log - warn", slog.LevelWarn)
	logging.PreSetupLog("Pre setup log - error", slog.LevelError)

	logging.Setup(true, output)

	slog.Info("Test logging", slog.String("test", "info"))
	slog.Debug("Test logging", slog.String("test", "debug"))
	slog.Warn("Test logging", slog.String("test", "warn"))
	slog.Error("Test logging", slog.String("test", "error"))

	assert.Contains(t, output.String(), "Pre setup log - debug")
	assert.Contains(t, output.String(), "Pre setup log - info")
	assert.Contains(t, output.String(), "Pre setup log - warn")
	assert.Contains(t, output.String(), "Pre setup log - error")
	assert.Contains(t, output.String(), "test=error")
	assert.Contains(t, output.String(), "test=debug")
	assert.Contains(t, output.String(), "test=info")
	assert.Contains(t, output.String(), "test=warn")
}
