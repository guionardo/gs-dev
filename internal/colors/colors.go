package colors

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
)

var (
	colorMap = map[string]*color.Color{
		NORMAL:    color.New(color.FgWhite).Add(color.BgBlack),
		PRIMARY:   color.New(color.FgBlue),
		SECONDARY: color.New(color.FgMagenta),
		TERTIARY:  color.New(color.FgRed),
		WARNING:   color.New(color.BgYellow).Add(color.FgBlack),
		ERROR:     color.New(color.BgRed).Add(color.FgWhite),
		SUCCESS:   color.New(color.FgGreen),
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
