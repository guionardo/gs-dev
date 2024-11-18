package calendar

import (
	"context"
	"testing"
	"time"
)

func TestRunProgress(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go RunProgress("Running", "Done", ctx)
		time.Sleep(6 * time.Second)
		cancel()
	})

}
