package url

func urlCommand(url string) (command string, args []string) {
	command = "rundll32"
	args = []string{"url.dll,FileProtocolHandler", url}
	return
}
