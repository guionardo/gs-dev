package calendar

const (
	Osc = "\u001B]"
	ST  = "\u001B\\"
)

func consoleUrl(url string, caption string) string {
	//https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
	// OSC + 8 + ; + ; +
	// printf '\e]8;;http://example.com\e\\This is a link\e]8;;\e\\\n'
	return Osc + "8;;" + url + ST + caption + Osc + "8;;" + ST
}
