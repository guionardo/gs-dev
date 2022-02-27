package dev

import (
	"fmt"
	"os"

	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/console"
	"github.com/spf13/cobra"
)

func DevGo(cmd *cobra.Command, args []string) error {
	cfg, err := configs.GetConfiguration(cmd)

	if err != nil {
		return err
	}
	output, _ := cmd.Flags().GetString("output")
	paths := cfg.DevConfig.Find(args)
	if len(paths) == 0 {
		console.OutputWarning("No paths for expression %v", args)
		return nil
	}

	if len(paths) == 1 {
		EchoGo(paths[0], output)
		return nil
	}
	items := make([]string, len(paths))
	for index, path := range paths {
		items[index] = path.FullPath
	}
	index, _, err := configs.PromptOpt("Path", items)
	if err != nil {
		return err
	}
	EchoGo(paths[index], output)
	return nil
}

func EchoGo(path configs.DevPathConfig, output string) {
	if len(output) == 0 {
		fmt.Printf("%v\n", path)
		return
	}
	lines := []string{
		fmt.Sprintf("cd %s", path.FullPath),
	}
	if len(path.AfterCommands) > 0 {
		lines = append(lines, path.AfterCommands...)
	}

	file, err := os.Create(output)
	if err != nil {
		fmt.Errorf("failed to create output file %s - %w", output, err)
		return
	}
	defer file.Close()
	for _, line := range lines {
		if _, err = fmt.Fprintf(file, "%s\n", line); err != nil {
			fmt.Errorf("failed to write to output file %s - %w\n", output, err)
			break
		}
	}
	if err != nil {
		os.Remove(output)
		fmt.Printf("Output file removed - %s", output)
	}

}
