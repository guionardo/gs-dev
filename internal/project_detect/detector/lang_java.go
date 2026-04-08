package detector

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Pom struct {
	XMLName      xml.Name `xml:"project"`
	ArtifactID   string   `xml:"artifactId"`
	Name         string   `xml:"name"`
	Description  string   `xml:"description"`
	Dependencies struct {
		Dependency []struct {
			GroupID    string `xml:"groupId"`
			ArtifactID string `xml:"artifactId"`
		} `xml:"dependency"`
	} `xml:"dependencies"`
}

func processPom(c *Candidate) {
	path := filepath.Join(c.Dir, "pom.xml")

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var p Pom
	if err := xml.Unmarshal(data, &p); err != nil {
		return
	}

	if p.Name != "" {
		c.Names = append(c.Names, p.Name)
	}

	if p.ArtifactID != "" {
		c.Names = append(c.Names, p.ArtifactID)
	}

	if s := strings.TrimSpace(p.Description); s != "" {
		c.addDesc(s, 0.95, "pom.xml")
	}

	for _, d := range p.Dependencies.Dependency {
		if d.GroupID != "" && d.ArtifactID != "" {
			c.Deps = append(c.Deps, normName(d.GroupID+"/"+d.ArtifactID))
			c.Deps = append(c.Deps, normName(d.ArtifactID))
		} else if d.ArtifactID != "" {
			c.Deps = append(c.Deps, normName(d.ArtifactID))
		}
	}

	c.Sources = append(c.Sources, "pom.xml")
}

func processGradle(c *Candidate) {
	path1 := filepath.Join(c.Dir, "build.gradle")
	path2 := filepath.Join(c.Dir, "build.gradle.kts")

	var path string
	if _, err := os.Stat(path1); err == nil {
		path = path1
	} else if _, err := os.Stat(path2); err == nil {
		path = path2
	} else {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	// description = "..."
	re := regexp.MustCompile(`(?m)^\s*description\s*=\s*["'](.+?)["']\s*$`)

	m := re.FindStringSubmatch(string(data))
	if len(m) == 2 {
		c.addDesc(strings.TrimSpace(m[1]), 0.9, filepath.Base(path))
	}
	// Simple detection of dependencies block lines like implementation("group:artifact")
	lines := strings.SplitSeq(string(data), "\n")
	for ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.Contains(ln, "implementation(") || strings.Contains(ln, "api(") {
			// try to extract artifact
			if i := strings.Index(ln, `"`); i >= 0 {
				if j := strings.Index(ln[i+1:], `"`); j >= 0 {
					ga := ln[i+1 : i+1+j]

					parts := strings.Split(ga, ":")
					if len(parts) >= 2 {
						c.Deps = append(c.Deps, normName(parts[1]))
						c.Deps = append(c.Deps, normName(parts[0]+"/"+parts[1]))
					}
				}
			}
		}
	}

	c.Sources = append(c.Sources, filepath.Base(path))
}
