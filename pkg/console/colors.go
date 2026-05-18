package console

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	NORMAL    = "normal"
	PRIMARY   = "primary"
	SECONDARY = "secondary"
	TERTIARY  = "tertiary"
	WARNING   = "warning"
	ERROR     = "error"
	SUCCESS   = "success"

	White  = "white"
	Cyan   = "cyan"
	Yellow = "yellow"
	Green  = "green"
	Red    = "red"
	Blue   = "blue"

	Bold = "bold"
)

var (
	colorMap = map[string]*color.Color{
		NORMAL:    color.New(color.FgWhite).Add(color.BgBlack),
		PRIMARY:   color.New(color.FgBlue).Add(color.FgBlack),
		SECONDARY: color.New(color.FgMagenta).Add(color.FgBlack),
		TERTIARY:  color.New(color.FgRed).Add(color.FgBlack),
		WARNING:   color.New(color.BgYellow).Add(color.FgBlack),
		ERROR:     color.New(color.BgRed).Add(color.FgWhite),
		SUCCESS:   color.New(color.FgGreen).Add(color.FgBlack),
	}
)

func Normal(msg string, a ...any) {
	_, _ = colorMap[NORMAL].Printf(msg, a...)
}

func Primary(msg string, a ...any) {
	_, _ = colorMap[PRIMARY].Printf(msg, a...)
}

func Secondary(msg string, a ...any) {
	_, _ = colorMap[SECONDARY].Printf(msg, a...)
}

func Tertiary(msg string, a ...any) {
	_, _ = colorMap[TERTIARY].Printf(msg, a...)
}

func Success(msg string, a ...any) {
	_, _ = colorMap[SUCCESS].Printf(msg, a...)
}

func Error(msg string, a ...any) {
	_, _ = colorMap[ERROR].Printf(msg, a...)
}

func Parse(msg string, a ...any) string {
	words := bytes.NewBufferString("")
	lastColor := NORMAL

	lastFnc := colorMap[lastColor].FprintFunc()
	for word := range strings.FieldsFuncSeq(msg, func(r rune) bool {
		return r == '{' || r == '}'
	}) {
		if color, ok := colorMap[word]; ok {
			lastFnc = color.FprintFunc()
			lastColor = word

			continue
		}

		lastFnc(words, word)
	}

	if lastColor != NORMAL {
		colorMap[NORMAL].FprintFunc()(words, "")
	}

	return fmt.Sprintf(words.String(), a...)
}

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
