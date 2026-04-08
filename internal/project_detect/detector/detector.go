package detector

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Languages   []string `json:"languages"`
	Frameworks  []string `json:"frameworks"`
	Path        string   `json:"path"`
	Confidence  float64  `json:"confidence"`
	Sources     []string `json:"sources"`
}

var ignoreDirs = map[string]struct{}{
	".git": {}, "node_modules": {}, "vendor": {}, ".venv": {}, "__pycache__": {},
	"dist": {}, "build": {}, "target": {}, "bin": {}, "obj": {}, ".tox": {},
	".idea": {}, ".vscode": {}, "coverage": {},
}

var codeExtToLang = map[string]string{
	".go": "go", ".py": "python", ".rb": "ruby", ".js": "javascript", ".ts": "typescript",
	".java": "java", ".kt": "kotlin", ".cs": "csharp", ".rs": "rust", ".php": "php",
	".swift": "swift", ".scala": "scala", ".c": "c", ".h": "c", ".hpp": "cpp",
	".cc": "cpp", ".cpp": "cpp", ".m": "objective-c", ".mm": "objective-c++",
	".tsx": "typescript", ".jsx": "javascript", ".sh": "shell", ".tf": "terraform",
}

var manifestNames = []string{
	packageJSON,
	pyprojectToml,
	goMod,
	cargoToml,
	composerJSON,
	"Chart.yaml",
	"serverless.yml",
	"serverless.yaml",
	"Pulumi.yaml",
	"Pulumi.yml",
	"pom.xml",
	"build.gradle",
	"build.gradle.kts",
	requirementsTxt,
	// csproj is pattern *.csproj
}

type ChartYAML struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ServerlessYAML struct {
	Service     any    `yaml:"service"`
	Description string `yaml:"description"`
}

type PulumiYAML struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func Detect(root string, maxFiles int) ([]Project, error) {
	projectDirs := discoverProjectDirs(root)

	projects := make([]Project, 0, len(projectDirs))
	for dir := range projectDirs {
		c := &Candidate{
			Dir:       dir,
			TypeHints: map[string]int{},
		}
		processManifests(c)
		// README
		if d, src, ok := readReadmeParagraph(dir); ok {
			c.addDesc(d, 0.8, src)
		}
		// Infer frameworks by deps
		c.Frameworks = detectFrameworks(c.Deps)
		// Infer type
		ptype := inferType(c)
		// Languages
		langs := detectLanguages(dir, maxFiles)
		name := chooseName(c, filepath.Base(dir))
		desc, conf, srcs := chooseDescription(c, ptype, langs, c.Frameworks)
		dedupSrc := dedupStrings(srcs)

		projects = append(projects, Project{
			Name:        name,
			Description: desc,
			Type:        ptype,
			Languages:   langs,
			Frameworks:  c.Frameworks,
			Path:        relPath(root, dir),
			Confidence:  conf,
			Sources:     dedupSrc,
		})
	}
	// Sort by path for stable ordering
	sort.Slice(projects, func(i, j int) bool { return projects[i].Path < projects[j].Path })

	return projects, nil
}

func relPath(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}

	if r == "." {
		return "."
	}

	return filepath.ToSlash(r)
}

func discoverProjectDirs(root string) map[string]struct{} {
	dirs := map[string]struct{}{}
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if d.IsDir() {
			base := d.Name()
			if _, ok := ignoreDirs[base]; ok && path != root {
				return filepath.SkipDir
			}

			return nil
		}

		base := filepath.Base(path)
		dir := filepath.Dir(path)

		// *.csproj
		if strings.HasSuffix(base, ".csproj") {
			dirs[dir] = struct{}{}
			return nil
		}
		// direct manifests
		if slices.Contains(manifestNames, base) {
			dirs[dir] = struct{}{}
			return nil
		}
		// pom.xml already covered by manifestNames
		return nil
	})

	return dirs
}

var fileSetProcessors = map[string]func(*Candidate){
	packageJSON:        processPackageJSON,
	pyprojectToml:      processPyProject,
	goMod:              processGoMod,
	cargoToml:          processCargoToml,
	composerJSON:       processComposerJSON,
	"Chart.yaml":       processChartYAML,
	"serverless.yml":   processServerless,
	"serverless.yaml":  processServerless,
	"Pulumi.yaml":      processPulumi,
	"Pulumi.yml":       processPulumi,
	requirementsTxt:    processRequirements,
	"csproj":           processCsprojIfAny,
	"pom.xml":          processPom,
	"build.gradle":     processGradle,
	"build.gradle.kts": processGradle,
}

func processManifests(c *Candidate) {
	fileSet := map[string]bool{}
	entries, _ := os.ReadDir(c.Dir)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fileSet[e.Name()] = true
	}

	for name, processor := range fileSetProcessors {
		if fileSet[name] {
			processor(c)
		}
	}
	// Terraform
	if hasTerraformFiles(c.Dir) {
		c.HasTerraform = true
		c.TypeHints["infra"]++
		c.Sources = append(c.Sources, "Terraform (*.tf)")
	}
	// OpenAPI
	if hasOpenAPI(c.Dir) {
		c.HasOpenAPI = true
		c.TypeHints["api"] += 2
		c.Sources = append(c.Sources, "OpenAPI spec")
	}
}

func processChartYAML(c *Candidate) {
	pathYaml := files.FindFirst(c.Dir, "Chart.yaml")

	ch, err := readYAMLFile[ChartYAML](pathYaml)
	if err != nil {
		return
	}

	if ch.Name != "" {
		c.Names = append(c.Names, ch.Name)
	}

	if strings.TrimSpace(ch.Description) != "" {
		c.addDesc(ch.Description, 1.0, "Chart.yaml")
	}

	c.TypeHints["infra"] += 2
	c.Sources = append(c.Sources, "Chart.yaml")
}

func processServerless(c *Candidate) {
	pathYaml := files.FindFirst(c.Dir, "serverless.yml", "serverless.yaml")

	sls, err := readYAMLFile[ServerlessYAML](pathYaml)
	if err != nil {
		return
	}
	// service may be string or object
	switch v := sls.Service.(type) {
	case string:
		if v != "" {
			c.Names = append(c.Names, v)
		}
	case map[string]any:
		if n, ok := v["name"].(string); ok && n != "" {
			c.Names = append(c.Names, n)
		}
	}

	if s := strings.TrimSpace(sls.Description); s != "" {
		c.addDesc(s, 0.95, "serverless.yml")
	}

	c.TypeHints["function"] += 2
	c.Sources = append(c.Sources, "serverless.yml")
}

func processPulumi(c *Candidate) {
	pathY := filepath.Join(c.Dir, "Pulumi.yaml")
	pathYml := filepath.Join(c.Dir, "Pulumi.yml")

	var path string
	if _, err := os.Stat(pathY); err == nil {
		path = pathY
	} else if _, err := os.Stat(pathYml); err == nil {
		path = pathYml
	} else {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var pl PulumiYAML
	if err := yaml.Unmarshal(data, &pl); err != nil {
		return
	}

	if pl.Name != "" {
		c.Names = append(c.Names, pl.Name)
	}

	if strings.TrimSpace(pl.Description) != "" {
		c.addDesc(pl.Description, 0.95, "Pulumi.yaml")
	}

	c.TypeHints["infra"]++
	c.Sources = append(c.Sources, "Pulumi.yaml")
}

func detectLanguages(dir string, maxFiles int) []string {
	counts := map[string]int{}
	files := 0
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if d.IsDir() {
			base := d.Name()
			if _, ok := ignoreDirs[base]; ok && path != dir {
				return filepath.SkipDir
			}

			return nil
		}

		if files > maxFiles {
			return io.EOF
		}

		files++

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if lang, ok := codeExtToLang[ext]; ok {
			counts[lang]++
		}

		return nil
	})

	type kv struct {
		K string
		V int
	}

	list := make([]kv, 0, len(counts))
	for k, v := range counts {
		list = append(list, kv{k, v})
	}

	sort.Slice(list, func(i, j int) bool { return list[i].V > list[j].V })

	out := []string{}
	for i := 0; i < len(list) && i < 3; i++ {
		out = append(out, list[i].K)
	}

	return out
}

func detectFrameworks(deps []string) []string {
	seen := map[string]struct{}{}
	add := func(s string) { seen[s] = struct{}{} }

	for _, d := range deps {
		// normalize
		switch {
		// Go
		case strings.Contains(d, "gin-gonic/gin"):
			add("gin")
		case strings.Contains(d, "labstack/echo"):
			add("echo")
		case strings.Contains(d, "gofiber/fiber"):
			add("fiber")
		case strings.Contains(d, "gorilla/mux"):
			add("gorilla/mux")
		case strings.Contains(d, "grpc"):
			add("grpc")
		case strings.Contains(d, "spf13/cobra"):
			add("cobra")
		case strings.Contains(d, "urfave/cli"):
			add("urfave/cli")
		// Python
		case d == "fastapi":
			add("fastapi")
		case d == "flask":
			add("flask")
		case d == "django":
			add("django")
		case d == "click":
			add("click")
		case d == "typer":
			add("typer")
		case d == "celery":
			add("celery")
		// Node
		case d == "express":
			add("express")
		case d == "fastify":
			add("fastify")
		case d == "koa":
			add("koa")
		case d == "nest" || d == "nestjs" || d == "@nestjs/core":
			add("nestjs")
		case d == "next" || d == "nextjs":
			add("nextjs")
		case d == "react" || d == "react-dom":
			add("react")
		case d == "vue":
			add("vue")
		case d == "angular" || d == "@angular/core":
			add("angular")
		case d == "svelte":
			add("svelte")
		// Java/Kotlin
		case strings.Contains(d, "spring-boot") || d == "spring-boot-starter-web":
			add("spring-boot")
		case strings.Contains(d, "quarkus"):
			add("quarkus")
		case strings.Contains(d, "micronaut"):
			add("micronaut")
		// Ruby
		case d == "rails" || d == "rubyonrails":
			add("rails")
		case d == "sinatra":
			add("sinatra")
		// Rust
		case d == "actix" || d == "actix-web":
			add("actix")
		case d == "axum":
			add("axum")
		case d == "tonic":
			add("tonic")
		case d == "clap":
			add("clap")
		// PHP
		case d == "laravel":
			add("laravel")
		case d == "symfony":
			add("symfony")
		// .NET (NuGet)
		case strings.HasPrefix(d, "microsoft.aspnetcore"):
			add("aspnetcore")
		case d == "serilog":
			add("serilog")
		case d == "hangfire":
			add("hangfire")
		case d == "masstransit":
			add("masstransit")
		}
	}
	// to slice
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func inferType(c *Candidate) string {
	score := map[string]int{}
	inc := func(t string, v int) { score[t] += v }

	// prior hints
	for k, v := range c.TypeHints {
		inc(k, v)
	}
	// Infra
	if c.HasTerraform {
		inc("infra", 3)
	}
	// Helm / Pulumi already count as infra

	// Frontend
	if contains(c.Frameworks, "react", "nextjs", "vue", "angular", "svelte") {
		inc("web-frontend", 3)
	}
	// API/back-end
	if contains(
		c.Frameworks,
		"gin",
		"echo",
		"fiber",
		"gorilla/mux",
		"grpc",
		"express",
		"fastify",
		"koa",
		"nestjs",
		"fastapi",
		"flask",
		"django",
		"spring-boot",
		"quarkus",
		"micronaut",
		"rails",
		"actix",
		"axum",
		"laravel",
		"symfony",
		"aspnetcore",
	) {
		inc("api", 3)
	}
	// CLI
	if contains(c.Frameworks, "cobra", "urfave/cli", "click", "typer", "clap") {
		inc("cli", 3)
	}
	// Function
	if v, ok := score["function"]; ok && v > 0 {
		// keep
	}
	// OpenAPI spec
	if c.HasOpenAPI {
		inc("api", 1)
	}

	// Decide
	type order struct {
		T string
		S int
	}

	var list []order
	for k, v := range score {
		list = append(list, order{k, v})
	}

	sort.Slice(list, func(i, j int) bool { return list[i].S > list[j].S })

	if len(list) > 0 && list[0].S > 0 {
		return list[0].T
	}
	// fallback
	return "application"
}

func chooseName(c *Candidate, fallback string) string {
	// prefer the most descriptive manifest name (no slash in final form)
	for _, n := range c.Names {
		if n == "" {
			continue
		}
		// if package path (e.g. github.com/org/proj), take last segment
		if strings.Contains(n, "/") {
			n = lastSegment(n)
		}
		// avoid leading/trailing whitespace noise
		n = strings.TrimSpace(n)
		if n != "" {
			return n
		}
	}

	return fallback
}

func chooseDescription(c *Candidate, ptype string, langs []string, frameworks []string) (string, float64, []string) {
	if len(c.DescCands) > 0 {
		// pick highest weight
		sort.Slice(c.DescCands, func(i, j int) bool { return c.DescCands[i].Weight > c.DescCands[j].Weight })
		best := c.DescCands[0]

		return best.Text, best.Weight, append([]string{}, c.Sources...)
	}
	// synthesize
	lang := ""
	if len(langs) > 0 {
		lang = langs[0]
	}

	var tech string

	if len(frameworks) > 0 {
		if len(frameworks) == 1 {
			tech = frameworks[0]
		} else {
			tech = frameworks[0] + ", " + frameworks[1]
		}
	}

	desc := synthesizeDescription(ptype, lang, tech, filepath.Base(c.Dir))

	return desc, 0.6, append([]string{}, c.Sources...)
}

func synthesizeDescription(ptype, lang, tech, fallbackName string) string {
	var t string

	switch ptype {
	case "api":
		t = "API"
	case "web-frontend":
		t = "web application (frontend)"
	case "cli":
		t = "command-line tool (CLI)"
	case "infra":
		t = "infrastructure module"
	case "function":
		t = "serverless function"
	default:
		t = "application"
	}

	var parts []string

	parts = append(parts, t)
	if lang != "" {
		parts = append(parts, "in "+cases.Title(language.English).String(lang))
	}

	if tech != "" {
		parts = append(parts, "using "+tech)
	}
	// domain fallback: common folder name patterns
	domain := inferDomainFromPath(fallbackName)
	if domain != "" {
		parts = append(parts, "for "+domain)
	}

	return strings.Join(parts, " ")
}

func inferDomainFromPath(name string) string {
	n := strings.ToLower(name)

	for _, kv := range []struct {
		K []string
		V string
	}{
		{[]string{"auth", "identity"}, "authentication"},
		{[]string{"order", "orders", "pedido"}, "order management"},
		{[]string{"billing", "invoice", "payment", "payments"}, "billing/payments"},
		{[]string{"etl", "ingest", "ingestion"}, "data ingestion/ETL"},
		{[]string{"admin", "backoffice"}, "administration"},
		{[]string{"web", "site", "frontend"}, "web interface"},
		{[]string{"worker", "consumer", "processor"}, "background processing"},
		{[]string{"report"}, "reporting"},
	} {
		for _, k := range kv.K {
			if strings.Contains(n, k) {
				return kv.V
			}
		}
	}

	return ""
}

func hasAny(m map[string]string, keys ...string) bool {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			return true
		}
	}

	return false
}

func getMap(m map[string]any, key string) (map[string]any, bool) {
	if v, ok := m[key]; ok {
		if mm, ok := v.(map[string]any); ok {
			return mm, true
		}
	}

	return nil, false
}

func getString(m map[string]any, key string) (string, bool) {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}

	return "", false
}

func normName(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "@") // npm scope

	return s
}

func lastSegment(s string) string {
	s = strings.TrimSuffix(s, "/")

	parts := strings.Split(s, "/")
	if len(parts) == 0 {
		return s
	}

	return parts[len(parts)-1]
}

func contains(slice []string, items ...string) bool {
	for item := range items {
		if slices.Contains(slice, items[item]) {
			return true
		}
	}

	return false
}

func dedupStrings(in []string) []string {
	set := map[string]struct{}{}
	out := []string{}

	for _, s := range in {
		if _, ok := set[s]; ok {
			continue
		}

		set[s] = struct{}{}
		out = append(out, s)
	}

	return out
}

func hasTerraformFiles(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if d.IsDir() {
			base := d.Name()
			if _, ok := ignoreDirs[base]; ok && path != dir {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(d.Name(), ".tf") {
			found = true
			return io.EOF
		}

		return nil
	})

	return found
}

func hasOpenAPI(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if d.IsDir() {
			base := d.Name()
			if _, ok := ignoreDirs[base]; ok && path != dir {
				return filepath.SkipDir
			}

			return nil
		}

		name := strings.ToLower(d.Name())
		if name == "openapi.yaml" || name == "openapi.yml" || name == "openapi.json" || name == "swagger.yaml" || name == "swagger.yml" || name == "swagger.json" {
			found = true
			return io.EOF
		}

		return nil
	})

	return found
}
