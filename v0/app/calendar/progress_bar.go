package calendar

import (
	"context"
	"fmt"
	"time"
)

func RunProgress(progress string, done string, ctx context.Context) {
	icons := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	iteration := 0
	// maxLength := max(len(progress), len(done))
	for {
		fmt.Printf("\r\x1B[2K%s %s", icons[iteration%len(icons)], progress)
		iteration++
		select {
		case <-ctx.Done():
			fmt.Printf("\r\x1B[2K⣿ %s : %s\n", progress, done)
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
