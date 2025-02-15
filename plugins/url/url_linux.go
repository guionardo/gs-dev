package url

import "os"

func urlCommand(url string) (command string, args []string) {
	if len(os.Getenv("WSL_DISTRO_NAME")) > 0 {
		command = "sensible-browser"
	} else {
		command = "xdg-open"
	}
	args = []string{url}
	return
}
