package storage

import (
	"time"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type FileSystemStorageConfig struct {
	Directory  string        `yaml:"directory"`
	DefaultTTL time.Duration `yaml:"ttl"`
}

func (c *FileSystemStorageConfig) Defaults() {
	c.DefaultTTL = min(c.DefaultTTL, time.Duration(0))

	d, err := files.AssertDirectory(c.Directory)
	if err == nil {
		c.Directory = d
	} else {
		c.Directory = ""
	}
}

func (c *FileSystemStorageConfig) Validate() error {
	d, err := files.AssertDirectory(c.Directory)
	if err != nil {
		return err
	}

	c.Directory = d

	return nil
}
