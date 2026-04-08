package detector

import "strings"

type (
	descCandidate struct {
		Text   string
		Weight float64
		Source string
	}

	Candidate struct {
		Dir          string
		Names        []string
		DescCands    []descCandidate
		Deps         []string
		Frameworks   []string
		TypeHints    map[string]int
		Sources      []string
		HasTerraform bool
		HasOpenAPI   bool
	}
)

func (c *Candidate) addDesc(text string, weight float64, src string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	// evita descrição idêntica repetida
	for _, d := range c.DescCands {
		if d.Text == text {
			return
		}
	}

	c.DescCands = append(c.DescCands, descCandidate{Text: text, Weight: weight, Source: src})
	c.Sources = append(c.Sources, src)
}
