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
	initScript = strings.ReplaceAll(initScript, "{GS_DEV}", executable)
	initScript = strings.ReplaceAll(initScript, "{GS_OUTPUT}", output)

	fmt.Print(initScript)
	return nil
}
