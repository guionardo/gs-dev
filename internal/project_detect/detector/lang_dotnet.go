package detector

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type CSProj struct {
	XMLName       xml.Name `xml:"Project"`
	PropertyGroup []struct {
		Description     string `xml:"Description"`
		AssemblyTitle   string `xml:"AssemblyTitle"`
		OutputType      string `xml:"OutputType"`
		TargetFramework string `xml:"TargetFramework"`
	} `xml:"PropertyGroup"`
	ItemGroup []struct {
		PackageReference []struct {
			Include string `xml:"Include,attr"`
			Version string `xml:"Version,attr"`
		} `xml:"PackageReference"`
	} `xml:"ItemGroup"`
}

func processCsprojIfAny(c *Candidate) {
	entries, _ := os.ReadDir(c.Dir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if strings.HasSuffix(e.Name(), ".csproj") {
			path := filepath.Join(c.Dir, e.Name())

			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			var pr CSProj
			if err := xml.Unmarshal(data, &pr); err != nil {
				continue
			}

			for _, pg := range pr.PropertyGroup {
				if s := strings.TrimSpace(pg.Description); s != "" {
					c.addDesc(s, 0.95, e.Name())
				}

				if s := strings.TrimSpace(pg.AssemblyTitle); s != "" {
					c.Names = append(c.Names, s)
				}

				if strings.Contains(strings.ToLower(pg.OutputType), "exe") {
					c.TypeHints["application"]++
				}
			}

			for _, ig := range pr.ItemGroup {
				for _, p := range ig.PackageReference {
					c.Deps = append(c.Deps, normName(p.Include))
				}
			}

			c.Names = append(c.Names, strings.TrimSuffix(e.Name(), ".csproj"))
			c.Sources = append(c.Sources, e.Name())
		}
	}
}
