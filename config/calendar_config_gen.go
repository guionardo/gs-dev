// CODE GENERATED. DO NOT EDIT.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func (cfg *CalendarsConfig) Save() (err error) {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()
	var content []byte
	if content, err = yaml.Marshal(cfg); err == nil {
		err = os.WriteFile(cfg.fileName, content, 0644)
	}
	return
}

func (cfg *CalendarsConfig) Load() (err error) {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()
	var content []byte
	if content, err = readFileWithLog(cfg.fileName); err == nil {
		err = yaml.Unmarshal(content, cfg)
	}
	return

}
