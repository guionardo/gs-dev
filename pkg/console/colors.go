package console

import (
	"strings"

	"github.com/fatih/color"
)

const (
	White  = "white"
	Cyan   = "cyan"
	Yellow = "yellow"
	Green  = "green"
	Red    = "red"
	Blue   = "blue"

	Bold = "bold"
)

// Styled format text adding color and style
func Styled(text, style string, args ...any) string {
	var c *color.Color

	switch {
	case strings.Contains(style, Cyan):
		c = color.New(color.FgCyan)
	case strings.Contains(style, Red):
		c = color.New(color.FgRed)
	case strings.Contains(style, Green):
		c = color.New(color.FgGreen)
	case strings.Contains(style, Yellow):
		c = color.New(color.FgYellow)
	case strings.Contains(style, Blue):
		c = color.New(color.FgBlue)
	default:
		c = color.New(color.FgWhite)
	}

	if strings.Contains(style, Bold) {
		c.Add(color.Bold)
	}

	return c.Sprintf(text, args...)
}
