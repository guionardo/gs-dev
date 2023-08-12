package url

func urlCommand(url string) (command string, args []string) {
	command = "open"
	args = []string{url}
	return
}
