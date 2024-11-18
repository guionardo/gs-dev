package pathtools

import (
	"os"
)

func CreatePath(path string) error {
	return os.Mkdir(path, os.ModeSticky|os.ModePerm)
}
