package todo

import "github.com/guionardo/gs-dev/config"

var readenConfig *config.TodoConfig

func getConfig() *config.TodoConfig {
	if readenConfig == nil {
		readenConfig = config.NewTodoConfig(config.GetConfigRepositoryFolder())
		_ = readenConfig.Load()
	}

	return readenConfig
}
