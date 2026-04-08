package detector

import (
	"path/filepath"
	"strings"
)

type ComposerJSON struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Require     map[string]string `json:"require"`
}

const composerJSON string = "composer.json"

func processComposerJSON(c *Candidate) {
	path := filepath.Join(c.Dir, composerJSON)

	comp, err := readJSONFile[ComposerJSON](path)
	if err != nil {
		return
	}

	if comp.Name != "" {
		c.Names = append(c.Names, comp.Name)
	}

	if s := strings.TrimSpace(comp.Description); s != "" {
		c.addDesc(s, 0.95, composerJSON)
	}

	for k := range comp.Require {
		c.Deps = append(c.Deps, normName(k))
	}

	c.Sources = append(c.Sources, composerJSON)
}
