package pathtools

import (
	"os"
	"syscall"
)

func CreatePath(path string) error {
	oldmask := syscall.Umask(0)
	defer syscall.Umask(oldmask)
	return os.Mkdir(path, os.ModeSticky|os.ModePerm)
}
