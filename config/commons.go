package config

import (
	"fmt"
	"os"
)

func readFileWithLog(filename string) (content []byte, err error) {
	if content, err = os.ReadFile(filename); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("[ERROR] Unexpected read error %v", err)
		}
	}
	return
}
