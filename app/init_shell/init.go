package initshell

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/app"
)

//go:embed init.sh
var initScript string

func InitAction() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	output := path.Join(os.TempDir(), app.ToolName)
	for key, value := range map[string]string{
		"GS_DEV":    executable,
		"GS_OUTPUT": output,
		"GS_TOOL":   app.ToolName + " v" + app.Version,
	} {
		initScript = strings.ReplaceAll(initScript, key, value)
	}

	fmt.Print(initScript)
	return nil
}
