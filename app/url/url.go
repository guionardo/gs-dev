package url

import (
	"fmt"
	"net/http"
	"os/exec"

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

func checkReachableUrl(url string) error {
	if resp, err := http.Head(url); err != nil {
		return err
	} else if resp.StatusCode >= 100 {
		return nil
	} else {
		return fmt.Errorf("got %s status from %s", resp.Status, url)
	}
}

func openInBrowser(url string) (err error) {
	if err = checkReachableUrl(url); err != nil {
		return
	}
	command, args := urlCommand(url)
	err = exec.Command(command, args...).Start()

	return
}
