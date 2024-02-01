package utils

import (
	"strings"
)

type PadType uint8

const (
	Left PadType = iota
	Right
	Center
)

func Pad(text string, length int, pad PadType) string {
	left, right := "", ""
	switch pad {
	case Left:
		if len(text) < length {
			right = strings.Repeat(" ", length-len(text))
		}

	case Right:
		if len(text) < length {
			left = strings.Repeat(" ", length-len(text))
		}
	case Center:
		if len(text) < length {
			left = strings.Repeat(" ", (length-len(text))/2)
			right = strings.Repeat(" ", length-len(text)-len(left))
		}
	}
	return left + text + right
}
