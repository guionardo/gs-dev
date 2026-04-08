package devservice

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	findpattern "github.com/guionardo/gs-dev/internal/find_pattern"
	"github.com/guionardo/gs-dev/internal/fs_tools"
	projectdetect "github.com/guionardo/gs-dev/internal/project_detect"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
)

type DevService struct {
	configFile  *config.ConfigFile
	devConfig   *DevConfiguration
	rootsConfig RootsConfiguration
}

func NewService(configuration *config.ConfigFile) *DevService {
	roots, err := config.GetValue[RootsConfiguration](configuration, "roots")
	if err != nil {
		slog.Warn("Error getting roots configuration", "error", err)
	}

	if roots == nil {
		roots = make(RootsConfiguration)
	}

	devConfig, err := config.GetValue[DevConfiguration](configuration, "dev")
	if err != nil {
		slog.Warn("Error getting dev configuration", "error", err)
	}

	devConfig.Defaults()

	return &DevService{
		configFile:  configuration,
		rootsConfig: roots,
		devConfig:   &devConfig,
	}
}

func (d *DevService) saveConfig() error {
	d.configFile.SetValue("dev", d.devConfig)
	d.configFile.SetValue("roots", d.rootsConfig)

	return d.configFile.Save()
}

func (d *DevService) AddRoot(root string) error {
	root, err := d.CanAddRoot(root)
	if err != nil {
		return err
	}

	d.rootsConfig[root] = Root{
		Folders:  []string{},
		MaxDepth: d.devConfig.DefaultMaxDepth,
	}

	return d.Sync()
}

// CanAddRoot checks if the root can be added and returns the absolute path
func (d *DevService) CanAddRoot(root string) (string, error) {
	root, err := fs_tools.AssertDirectory(root)
	if err != nil {
		return root, err
	}

	if _, ok := d.rootsConfig[root]; ok {
		return root, fmt.Errorf("root %s already exists", root)
	}

	return root, nil
}

func (d *DevService) DeleteRoot(root string) error {
	if _, ok := d.rootsConfig[root]; !ok {
		return fmt.Errorf("root %s not found", root)
	}

	delete(d.rootsConfig, root)

	return d.saveConfig()
}

func (d *DevService) CanDeleteRoot(root string) bool {
	_, ok := d.rootsConfig[root]
	return ok
}

func (d *DevService) Sync() error {
	roots := make(RootsConfiguration)

	colors.Primary("Syncing %d roots\n", len(d.rootsConfig))

	for directory, root := range d.rootsConfig {
		err := syncRoot(directory, &root)
		if err != nil {
			return err
		}

		roots[directory] = root
	}

	d.devConfig.LastSync = time.Now()
	d.rootsConfig = roots

	return d.saveConfig()
}

func (d *DevService) IsTimeToSync() bool {
	return d.devConfig.LastSync.Add(d.devConfig.SyncInterval).Before(time.Now())
}

func (d *DevService) DoSyncIfNeeded() error {
	if !d.IsTimeToSync() {
		slog.Debug("Not time to sync", slog.String("last_sync", d.devConfig.LastSync.Format(time.RFC3339)))
		return nil
	}

	slog.Debug("Time to sync", slog.String("last_sync", d.devConfig.LastSync.Format(time.RFC3339)))

	return d.Sync()
}

func (d *DevService) GetRoots() ([]string, error) {
	roots := make([]string, 0, len(d.rootsConfig))
	for root := range d.rootsConfig {
		roots = append(roots, root)
	}

	return roots, nil
}

func (d *DevService) ListRoots() error {
	roots, err := d.GetRoots()
	if err != nil {
		return err
	}

	for _, root := range roots {
		colors.Primary("\nRoot: %s\n", root)

		for _, folder := range d.rootsConfig[root].Folders {
			project, err := projectdetect.DetectProject(folder)
			if err == nil {
				colors.Secondary("\t %s\n", project.String())
			}
		}
	}

	return nil
}

func (d *DevService) PurgeUnexistentRoots() (err error) {
	deleted := false

	for root := range d.rootsConfig {
		if !pathtools.DirExists(root) {
			slog.Info("Purging unexistent root", slog.String("directory", root))
			delete(d.rootsConfig, root)

			deleted = true
		}
	}

	if deleted {
		err = d.saveConfig()

		slog.Info("Roots purged")
	}

	return err
}

// RunFind finds the folders that match the words
func (d *DevService) RunFind(words []string) (err error) {
	if err = d.DoSyncIfNeeded(); err != nil {
		return err
	}

	folders := d.GetFilteredFolders(words)
	if len(folders) == 0 {
		return fmt.Errorf("no folders found for %v", words)
	}

	folder, err := dialog.Choose("Choose a folder:", (dialog.ToAnyArray(folders))...)
	if err == nil {
		d.ChosenFolder(folder)
		postcommand.AddOutputLine("cd " + folder)
	}

	return err
}

func (d *DevService) ChosenFolder(folder string) {
	d.devConfig.LastChosen[folder]++
	_ = d.saveConfig()
}

func (d *DevService) GetFilteredFolders(words []string) (folders []string) {
	for rootFolder, root := range d.rootsConfig {
		for index := range root.Folders {
			if len(words) == 0 || findpattern.PathContainsPattern(root.Folders[index], rootFolder, words) {
				folders = append(folders, root.Folders[index])
			}
		}
	}

	return folders
}

func (d *DevService) RunFavorites() (err error) {
	if len(d.devConfig.LastChosen) == 0 {
		return errors.New("no folders chosen")
	}

	lastChosen := make([]string, 0, len(d.devConfig.LastChosen))
	for folder, count := range d.devConfig.LastChosen {
		lastChosen = append(lastChosen, fmt.Sprintf("%06d%s", count, folder))
	}
	// Sort by count descending
	sort.Slice(lastChosen, func(i, j int) bool {
		return lastChosen[i] > lastChosen[j]
	})

	folders := make([]string, 0, len(lastChosen))
	for i := range min(len(lastChosen), d.devConfig.MostChosenCount) {
		folders = append(folders, lastChosen[i][6:])
	}

	folder, err := dialog.Choose("Choose a folder:", (dialog.ToAnyArray(folders))...)
	if err == nil {
		d.ChosenFolder(folder)
		postcommand.AddOutputLine("cd " + folder)
	}

	return nil
}

func (d *DevService) Setup() error {
	return nil
}

func syncRoot(directory string, root *Root) error {
	reader := NewRootReader(directory, root)
	if root.LocalConfigs == nil {
		root.LocalConfigs = make(map[string]LocalConfig)
	}

	removed, added := reader.SyncSummary(root.CanIncludeFolder)
	defer root.Resync()

	if len(removed) == 0 && len(added) == 0 {
		colors.Normal("Root [%s] is up to date\n", directory)
		return nil
	}

	reader.UpdateRoot(root)
	colors.Primary("Root [%s] has %d changes:\n", directory, len(removed)+len(added))

	for _, project := range removed {
		colors.Secondary("\t%s removed\n", project.Folder)
		delete(root.LocalConfigs, project.Folder)
	}

	for _, project := range added {
		colors.Success("\t%s added\n", project.String())

		if _, ok := root.LocalConfigs[project.Folder]; !ok {
			localConfig, err := NewLocalConfig(project.Folder)
			if err == nil {
				root.LocalConfigs[project.Folder] = *localConfig
			}
		}
	}

	return nil
}
