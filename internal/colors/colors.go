package colors

import (
	"github.com/fatih/color"
)

var (
	normal = color.New(color.FgWhite).Printf
	red    = color.New(color.FgRed).Printf
	green  = color.New(color.FgGreen).Printf
	blue   = color.New(color.FgBlue).Printf
)

func Normal(msg string, a ...interface{}) {
	_, _ = normal(msg, a...)
}

func Red(msg string, a ...interface{}) {
	_, _ = red(msg, a...)
}

func Green(msg string, a ...interface{}) {
	_, _ = green(msg, a...)
}

func Blue(msg string, a ...interface{}) {
	_, _ = blue(msg, a...)
}
