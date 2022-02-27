package dev

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/console"
	"github.com/spf13/cobra"
)

func DevSetup(cmd *cobra.Command) (err error) {
	status, _ := cmd.Flags().GetString("status")
	var enable, disable, changed bool
	switch strings.ToLower(status) {
	case "enabled":
		enable = true
	case "disabled":
		disable = true
	}

	var cfg *configs.RootConfig
	if cfg, err = configs.GetConfiguration(cmd); err != nil {
		return
	}
	var pwd string
	if pwd, err = os.Getwd(); err != nil {
		return
	}
	var index int
	var path *configs.DevPathConfig
	if index, path = cfg.DevConfig.FindPath(pwd); index == -1 {
		return fmt.Errorf("current path is not a dev folder: %s", pwd)
	}

	if enable && path.IsHidden {
		path.IsHidden = false
		changed = true
	}
	if disable && !path.IsHidden {
		path.IsHidden = true
		changed = true
	}

	afterCommands, _ := cmd.Flags().GetStringArray("after-command")
	if !reflect.DeepEqual(path.AfterCommands, afterCommands) {
		path.AfterCommands = afterCommands
		changed = true
	}

	ignore, err := cmd.Flags().GetBool("ignore-subpaths")
	if path.IgnoreSubPaths != ignore {
		path.IgnoreSubPaths = ignore
		doIgnoreSubPaths(cfg, path.FullPath, ignore)
		changed = true
	}

	if changed {
		if err = cfg.WriteFile(); err == nil {
			console.OutputSuccess("Updated folder %v", path)
		}
	} else {
		console.OutputNeutral("No changes to folder %v", path)
	}
	return
}

func doIgnoreSubPaths(cfg *configs.RootConfig, basePath string, ignore bool) int {
	changed := 0
	for _, path := range cfg.DevConfig.Paths {
		if len(path.FullPath) > len(basePath) && strings.HasPrefix(path.FullPath, basePath) {
			path.IsHidden = ignore
			changed++
		}
	}
	return changed
}
