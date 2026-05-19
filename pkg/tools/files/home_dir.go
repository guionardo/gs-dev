package files

import "os"

var homeDir string

func init() {
	homeDir, _ = os.UserHomeDir()
}
