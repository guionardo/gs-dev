package detector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

const (
	pyprojectToml   string = "pyproject.toml"
	requirementsTxt string = "requirements.txt"
)

func processPyProject(c *Candidate) {
	path := files.FindFirst(c.Dir, pyprojectToml)

	tree, err := readTOMLFile(path)
	if err != nil {
		return
	}

	// [project]
	if proj, ok := getMap(tree, "project"); ok {
		if name, ok := getString(proj, "name"); ok {
			c.Names = append(c.Names, name)
		}

		if desc, ok := getString(proj, "description"); ok && strings.TrimSpace(desc) != "" {
			c.addDesc(desc, 0.95, "pyproject.toml [project]")
		}
		// dependencies (array)
		if deps, ok := getSlice(proj, "dependencies"); ok {
			for _, v := range deps {
				if s, ok := v.(string); ok {
					c.Deps = append(c.Deps, normReqName(s))
				}
			}
		}
	}
	// [tool.poetry]
	if tool, ok := getMap(tree, "tool"); ok {
		if poetry, ok := getMap(tool, "poetry"); ok {
			if name, ok := getString(poetry, "name"); ok {
				c.Names = append(c.Names, name)
			}

			if desc, ok := getString(poetry, "description"); ok && strings.TrimSpace(desc) != "" {
				c.addDesc(desc, 0.9, "pyproject.toml [tool.poetry]")
			}
			// dependencies (map)
			if deps, ok := getMap(poetry, "dependencies"); ok {
				for k := range deps {
					if strings.ToLower(k) == "python" {
						continue
					}

					c.Deps = append(c.Deps, normName(k))
				}
			}
		}
	}

	c.Sources = append(c.Sources, pyprojectToml)
}

func processRequirements(c *Candidate) {
	path := filepath.Join(c.Dir, requirementsTxt)

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.SplitSeq(string(data), "\n")
	for ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}

		c.Deps = append(c.Deps, normReqName(ln))
	}

	c.Sources = append(c.Sources, requirementsTxt)
}
func normReqName(s string) string {
	// take package name before version spec and extras
	s = strings.TrimSpace(strings.ToLower(s))
	if i := strings.IndexAny(s, " <>=!"); i >= 0 {
		s = s[:i]
	}
	// remove extras like fastapi[all]
	if i := strings.IndexByte(s, '['); i >= 0 {
		s = s[:i]
	}
	// for URLs or git+https, take the last segment
	if strings.HasPrefix(s, "git+") || strings.HasPrefix(s, "http") {
		s = lastSegment(s)
	}
	// for "pkg==1.2.3"
	if i := strings.Index(s, "=="); i >= 0 {
		s = s[:i]
	}

	return strings.TrimSpace(s)
}
func getSlice(m map[string]any, key string) ([]any, bool) {
	if v, ok := m[key]; ok {
		if s, ok := v.([]any); ok {
			return s, true
		}
	}

	return nil, false
}
