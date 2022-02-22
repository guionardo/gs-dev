package configs

type DevConfig struct {
	DevFolders   []string `yaml:"dev_folders"`
	MaxSubLevels int      `yaml:"max_sub_levels"`
}

func (cfg *DevConfig) Prompt() bool {
	if !Confirm("Setup dev config?", true) {
		return false
	}
	return true
}
