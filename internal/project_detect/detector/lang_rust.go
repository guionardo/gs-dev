package detector

import (
	"path/filepath"
	"strings"
)

const cargoToml string = "Cargo.toml"

func processCargoToml(c *Candidate) {
	path := filepath.Join(c.Dir, cargoToml)

	tree, err := readTOMLFile(path)
	if err != nil {
		return
	}

	if pkg, ok := getMap(tree, "package"); ok {
		if name, ok := getString(pkg, "name"); ok {
			c.Names = append(c.Names, name)
		}

		if desc, ok := getString(pkg, "description"); ok && strings.TrimSpace(desc) != "" {
			c.addDesc(desc, 0.95, cargoToml)
		}
	}
	// [dependencies]
	if deps, ok := getMap(tree, "dependencies"); ok {
		for k := range deps {
			c.Deps = append(c.Deps, normName(k))
		}
	}

	c.Sources = append(c.Sources, cargoToml)
}
