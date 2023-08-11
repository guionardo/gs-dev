package url

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/guionardo/gs-dev/internal/git"
)

func RunUrl(pathName string, justShow bool) (err error) {
	var url string
	if url, err = git.GetRemoteHttpURL(pathName); err == nil {
		fmt.Print(url)
		if !justShow {
			err = openInBrowser(url)
		}
	}
	return
}

func openInBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		command := "xdg-open"
		if len(os.Getenv("WSL_DISTRO_NAME")) > 0 {
			command = "sensible-browser"
		}
		err = exec.Command(command, url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return
}
