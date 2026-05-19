package postcommand

import (
	"fmt"
	"os"
)

const outputFilePermissions = 0600

var (
	outputContent         []byte
	postCommandOutputFile string
)

func SetOutputFile(outputFile string) {
	postCommandOutputFile = outputFile
}

func AddOutputLine(line string) {
	outputContent = fmt.Appendf(outputContent, "%s\n", line)
}

func WriteOutput() error {
	if len(postCommandOutputFile) == 0 {
		return nil
	}

	err := os.WriteFile(postCommandOutputFile, outputContent, outputFilePermissions)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	outputContent = []byte{}

	return nil
}
