package logging

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Parallel()

	output := bytes.NewBufferString("")

	_ = Logger("--debug", output)

	Debug("Pre setup log - debug")
	Info("Pre setup log - info")
	Warn("Pre setup log - warn")
	Error("Pre setup log - error")

	Info("Test logging", slog.String("test", "info"))
	Debug("Test logging", slog.String("test", "debug"))
	Warn("Test logging", slog.String("test", "warn"))
	Error("Test logging", slog.String("test", "error"))

	assert.Contains(t, output.String(), "Pre setup log - debug")
	assert.Contains(t, output.String(), "Pre setup log - info")
	assert.Contains(t, output.String(), "Pre setup log - warn")
	assert.Contains(t, output.String(), "Pre setup log - error")
	assert.Contains(t, output.String(), "test=error")
	assert.Contains(t, output.String(), "test=debug")
	assert.Contains(t, output.String(), "test=info")
	assert.Contains(t, output.String(), "test=warn")
}
