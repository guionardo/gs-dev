package detector

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var readmeNames = []string{"README.md", "README.MD", "Readme.md", "README", "README.rst", "README.txt"}

func readReadmeParagraph(dir string) (string, string, bool) {
	for _, n := range readmeNames {
		p := filepath.Join(dir, n)

		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}

		text := firstMeaningfulParagraph(string(data))
		if text != "" {
			return text, n, true
		}
	}

	return "", "", false
}

func firstMeaningfulParagraph(md string) string {
	// Remove badges (linhas com [![ ou <img)
	lines := strings.Split(md, "\n")

	var buf bytes.Buffer

	par := []string{}
	flush := func() string {
		t := strings.TrimSpace(strings.Join(par, " "))
		t = stripMarkdown(t)
		// filtros simples
		if t == "" || len(t) < 10 {
			return ""
		}

		return t
	}

	for _, ln := range lines {
		s := strings.TrimSpace(ln)
		if s == "" {
			if len(par) > 0 {
				if t := flush(); t != "" {
					return t
				}

				par = par[:0]
			}

			continue
		}
		// skip headings and badges
		if strings.HasPrefix(s, "#") || strings.HasPrefix(s, "<h") {
			continue
		}

		if strings.Contains(s, "[![") || strings.Contains(s, "<img") {
			continue
		}

		par = append(par, s)

		if buf.Len() > 2000 {
			break
		}
	}

	if len(par) > 0 {
		if t := flush(); t != "" {
			return t
		}
	}

	return ""
}

var (
	badgesRegex  = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`!\[[^\]]*\]\([^\)]*\)`) })
	badgesRegex2 = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`\[[^\]]*\]\([^\)]*\)`) })
)

func stripMarkdown(s string) string {
	// remove links [text](url), inline code, bold/italics
	s = badgesRegex().ReplaceAllString(s, "")
	s = badgesRegex2().ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.TrimSpace(s)

	return s
}
