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
	var cfg *configs.RootConfig
	if cfg, err = configs.GetConfiguration(cmd); err != nil {
		return
	}
	_, path, err := getCurrentPath(cfg)
	if err != nil {
		return err
	}
	changed := parseEnableDisable(path, cmd)
	changed = changed || parseAfterCommand(path, cmd)
	changed = changed || parseIgnoreSubPaths(path, cmd, cfg)

	if changed {
		if err = cfg.WriteFile(); err == nil {
			console.OutputSuccess("Updated folder %v", path)
		}
	} else {
		console.OutputNeutral("No changes to folder %v", path)
	}
	return
}

func getCurrentPath(cfg *configs.RootConfig) (pathIndex int, path *configs.DevPathConfig, err error) {
	var pwd string
	if pwd, err = os.Getwd(); err != nil {
		return -1, nil, err
	}
	if pathIndex, path = cfg.DevConfig.FindPath(pwd); pathIndex == -1 {
		return -1, nil, fmt.Errorf("current path is not a dev folder: %s", pwd)
	}
	return
}

func parseEnableDisable(path *configs.DevPathConfig, cmd *cobra.Command) (changed bool) {
	status, err := cmd.Flags().GetString("status")
	changed = false
	if err != nil {
		return
	}

	switch strings.ToLower(status) {
	case "enabled":
		if path.IsHidden {
			path.IsHidden = false
			changed = true
		}
	case "disabled":
		if !path.IsHidden {
			path.IsHidden = true
			changed = true
		}
	}
	return
}

func parseAfterCommand(path *configs.DevPathConfig, cmd *cobra.Command) (changed bool) {
	changed = false
	afterCommands, err := cmd.Flags().GetStringArray("after-command")
	if err == nil && !reflect.DeepEqual(path.AfterCommands, afterCommands) {
		path.AfterCommands = afterCommands
		changed = true
	}
	return
}

func parseIgnoreSubPaths(path *configs.DevPathConfig, cmd *cobra.Command, cfg *configs.RootConfig) (changed bool) {
	changed = false
	if ignore, err := cmd.Flags().GetBool("ignore-subpaths"); err == nil && path.IgnoreSubPaths != ignore {
		path.IgnoreSubPaths = ignore
		doIgnoreSubPaths(cfg, path.FullPath, ignore)
		changed = true
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
