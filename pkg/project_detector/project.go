package projectdetector

import (
	"fmt"

	"github.com/guionardo/gs-dev/pkg/console"
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
		style = console.White
	}

	return fmt.Sprintf("%s %s: %s - %s", console.Styled(pd.Symbol(), style), console.Styled(pd.Type, style+" bold"), pd.Name, pd.Folder)
}
