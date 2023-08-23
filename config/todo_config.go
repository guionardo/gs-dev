package config

import (
	"path"
	"sort"
	"time"
)

type TodoConfig struct {
	fileName        string
	Todos           []*TodoItem   `yaml:"todos"`
	LastSync        time.Time     `yaml:"last_sync"`
	MaxSyncInterval time.Duration `yaml:"max_sync_interval"`
}

type TodoItem struct {
	Name      string    `yaml:"name"`
	DueTo     time.Time `yaml:"due_to"`
	Completed bool      `yaml:"completed"`
	Tags      []string  `yaml:"tags"`
}

func NewTodoConfig(root string) (cfg *TodoConfig) {
	return &TodoConfig{
		fileName:        path.Join(root, "todo.yaml"),
		Todos:           make([]*TodoItem, 0),
		MaxSyncInterval: time.Hour,
	}
}

func (cfg *TodoConfig) GetAllTags() []string {
	tags := make(map[string]struct{}, len(cfg.Todos))
	for _, todo := range cfg.Todos {
		for _, tag := range todo.Tags {
			tags[tag] = struct{}{}
		}
	}
	tagList := make([]string, len(tags))
	index := 0
	for tag := range tags {
		tagList[index] = tag
		index++
	}
	return sort.StringSlice(tagList)
}
