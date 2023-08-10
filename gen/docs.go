//go:build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/guionardo/gs-dev/cmd"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := cmd.GetRootCmd()
	docsFolder, err := filepath.Abs("../docs")
	if err != nil {
		log.Fatalf("Failed to get ../docs folder - %v", err)
	}
	if err = os.RemoveAll(docsFolder); err != nil {
		log.Fatalf("Failed to remove folder %s - %v", docsFolder, err)
	}
	if err = pathtools.CreatePath(docsFolder); err != nil {
		log.Fatal("Failed to create folder %s - %v", docsFolder, err)
	}
	if err = doc.GenMarkdownTree(rootCmd, "../docs"); err != nil {
		log.Fatalf("Failed to generate documentation in %s - %v", docsFolder, err)
	}
	log.Printf("Documentation created at %s", docsFolder)
}
