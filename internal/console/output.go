package console

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type OutputLevel int8

const (
	INFO OutputLevel = iota
	WARN
	SUCCESS
	ERROR
	NONE
)

type iconStruct struct {
	icon         string
	messageColor func(interface{}) string
}

var icons map[OutputLevel]iconStruct

func init() {
	icons = map[OutputLevel]iconStruct{
		INFO:    {promptui.Styler(promptui.FGBlue)("ℹ️"), promptui.Styler(promptui.BGBlue, promptui.FGWhite)},
		WARN:    {promptui.IconWarn, promptui.Styler(promptui.BGYellow, promptui.FGBlack)},
		SUCCESS: {promptui.IconGood, promptui.Styler(promptui.BGGreen, promptui.FGWhite)},
		ERROR:   {promptui.IconBad, promptui.Styler(promptui.BGRed, promptui.FGWhite)},
		NONE:    {"", promptui.Styler()},
	}
}

func Output(printLevel OutputLevel, format string, args ...interface{}) {
	fmt.Printf("%s %s\n", icons[printLevel].icon, icons[printLevel].messageColor(fmt.Sprintf(format, args...)))
}

func OutputInfo(format string, args ...interface{}) {
	Output(INFO, format, args...)
}

func OutputError(format string, args ...interface{}) {
	Output(ERROR, format, args...)
}

func OutputSuccess(format string, args ...interface{}) {
	Output(SUCCESS, format, args...)
}

func OutputWarning(format string, args ...interface{}) {
	Output(WARN, format, args...)
}

func OutputNeutral(format string, args ...interface{}) {
	Output(NONE, format, args...)
}
