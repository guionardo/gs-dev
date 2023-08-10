package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var templ = `// CODE GENERATED. DO NOT EDIT.
package config

import (
	"os"	

	"gopkg.in/yaml.v3"
)

func (cfg *{{STRUCT}}) Save() error {
	if content, err := yaml.Marshal(cfg); err != nil {
		return err
	} else {
		return os.WriteFile(cfg.fileName, content, 0644)
	}
}

func (cfg *{{STRUCT}}) Load() error {
	if content, err := os.ReadFile(cfg.fileName); err != nil {
		return err
	} else {
		return yaml.Unmarshal(content, cfg)
	}
}
`

func generate(filePath string) error {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`type\s(?P<struct>[A-Za-z]{1,20})\sstruct\s{`)
	matches := re.FindStringSubmatch(string(source))
	structIndex := re.SubexpIndex("struct")
	structName := matches[structIndex]

	if len(structName) == 0 {
		return fmt.Errorf("struct not found in file %s", filePath)
	}
	log.Printf("File %s -> %s\n", filePath, structName)

	content := strings.ReplaceAll(templ, "{{STRUCT}}", structName)
	genFilePath := strings.TrimSuffix(filePath, ".go") + "_gen.go"
	return os.WriteFile(genFilePath, []byte(content), 0644)
}

func main() {
	configPath, _ := filepath.Abs(".")

	log.Printf("Generate %s\n", configPath)
	files, err := os.ReadDir(configPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_config.go") {
			fileName := path.Join(configPath, file.Name())
			if err = generate(fileName); err != nil {
				log.Fatal(err)
			}
		}
	}
}
