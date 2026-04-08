package detector

import (
	"os"
	"path/filepath"
	"strings"
)

const goMod string = "go.mod"

func processGoMod(c *Candidate) {
	path := filepath.Join(c.Dir, goMod)

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.SplitSeq(string(data), "\n")
	for ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "module ") {
			mod := strings.TrimSpace(strings.TrimPrefix(ln, "module"))

			mod = strings.TrimSpace(mod)
			if mod != "" {
				c.Names = append(c.Names, lastSegment(mod))
			}
		}

		if strings.HasPrefix(ln, "require ") || (!strings.HasPrefix(ln, "replace") && strings.Contains(ln, " ")) {
			// try to extract module names in simple way
			parts := strings.Fields(ln)
			// typical: require github.com/gin-gonic/gin v1.9.0
			if len(parts) >= 2 && strings.Contains(parts[0], "require") {
				if len(parts) >= 3 {
					c.Deps = append(c.Deps, normName(parts[1]))
				}
			} else if len(parts) >= 2 && isGoModulePath(parts[0]) {
				c.Deps = append(c.Deps, normName(parts[0]))
			}
		}
	}

	c.Sources = append(c.Sources, goMod)
}
func isGoModulePath(s string) bool {
	return strings.Count(s, ".") >= 1 && strings.Contains(s, "/")
}
