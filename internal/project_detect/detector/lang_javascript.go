package detector

import (
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

const (
	packageJSON string = "package.json"
)

type (
	PackageJSON struct {
		Name            string            `json:"name"`
		Description     string            `json:"description"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Scripts         map[string]string `json:"scripts"`
	}
)

var jsDependencies = []string{
	"next", "react", "vue", "angular", "nuxt", "svelte",
}

func processPackageJSON(c *Candidate) {
	path := files.FindFirst(c.Dir, packageJSON)

	pkg, err := readJSONFile[PackageJSON](path)
	if err != nil {
		return
	}

	if pkg.Name != "" {
		c.Names = append(c.Names, pkg.Name)
	}

	if s := strings.TrimSpace(pkg.Description); s != "" {
		c.addDesc(s, 0.95, packageJSON)
	}

	for k := range pkg.Dependencies {
		c.Deps = append(c.Deps, normName(k))
	}

	for k := range pkg.DevDependencies {
		c.Deps = append(c.Deps, normName(k))
	}
	// Front hints
	if hasAny(pkg.Dependencies, jsDependencies...) ||
		hasAny(pkg.DevDependencies, jsDependencies...) {
		c.TypeHints["web-frontend"] += 2
	}

	c.Sources = append(c.Sources, packageJSON)
}
