package config

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/guionardo/gs-dev/app/build"
)

const CONFIG_DIR_ENV = "GS_DEV_CONFIG_DIR"

func getConfigDir(appName string) (configDir string, err error) {
	if configDir = os.Getenv(CONFIG_DIR_ENV); len(configDir) == 0 {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		configDir = path.Join(user.HomeDir, ".config", appName)
	}
	stat, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
	} else if !stat.IsDir() {
		err = fmt.Errorf("config dir is not a directory: %s", configDir)
	}
	return
}

func GetConfigDir(appName ...string) string {
	if len(appName) == 0 {
		appName = []string{build.AppName}
	}
	return sync.OnceValue(func() string {
		if cd, err := getConfigDir(appName[0]); err != nil {
			panic(fmt.Sprintf("error getting config dir: %s - %v", cd, err))
		} else {
			return cd
		}
	})()
}
