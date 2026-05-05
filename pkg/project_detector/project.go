package projectdetector

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type (
	ProjectData struct {
		Folder      string // Absolute path to the project folder
		Type        string // Detected project type, e.g., "node", "python", "dotnet", etc.
		Name        string // Detected project name, e.g., from package.json, .csproj, etc.
		Description string // TODO: Implement parser for reading description from source
	}
)

func (pd ProjectData) Symbol() string {
	if symbol, ok := symbols[pd.Type]; ok {
		return symbol
	}

	return symbols["unknown"]
}
func (pd ProjectData) String() string {
	return fmt.Sprintf("%s %s: %s - %s", pd.Symbol(), pd.Type, pd.Name, pd.Folder)
}

func (pd ProjectData) ColoredString() string {
	style, ok := styles[pd.Type]
	if !ok {
		style = White
	}

	return fmt.Sprintf("%s %s: %s - %s", styled(pd.Symbol(), style), styled(pd.Type, style+" bold"), pd.Name, pd.Folder)
}

func styled(text, style string) string {
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

	if strings.Contains(style, "bold") {
		c.Add(color.Bold)
	}

	return c.Sprint(text)
}
