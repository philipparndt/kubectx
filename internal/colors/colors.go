package colors

import "github.com/fatih/color"

var Red = color.New(color.FgRed).SprintFunc()
var Cyan = color.New(color.FgCyan).SprintFunc()
var Gray = color.New(color.FgHiBlack).SprintFunc()

var ColorGray = "\033[90m"
var ColorReset = "\033[0m"

func init() {
	if color.NoColor {
		ColorGray = ""
		ColorReset = ""
	}
}
