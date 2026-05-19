package dialog

import (
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
)

type ConfigDialog struct {
	SyncInterval    time.Duration
	DefaultMaxDepth int
	MostChosenCount int
}

func ReadConfig(syncInterval time.Duration, defaultMaxDepth int, mostChosenCount int) (*ConfigDialog, error) {
	config := &ConfigDialog{
		SyncInterval:    syncInterval,
		DefaultMaxDepth: defaultMaxDepth,
		MostChosenCount: mostChosenCount,
	}
	values := map[string]any{
		"syncInterval":    config.SyncInterval.String(),
		"defaultMaxDepth": strconv.Itoa(config.DefaultMaxDepth),
		"mostChosenCount": strconv.Itoa(config.MostChosenCount),
	}
	intervals := []time.Duration{
		15 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		2 * time.Hour,
		3 * time.Hour,
		4 * time.Hour,
		5 * time.Hour,
		6 * time.Hour,
		7 * time.Hour,
		8 * time.Hour,
		9 * time.Hour,
		10 * time.Hour,
	}

	intervalsValues := make([]string, len(intervals))
	for i, interval := range intervals {
		intervalsValues[i] = interval.String()
	}

	questions := []*survey.Question{
		{
			Name: "syncInterval",
			Prompt: &survey.Select{
				Message: "Sync interval",
				Options: intervalsValues,
				Default: config.SyncInterval.String(),
				Help:    "The interval between reading the roots for synchronization",
			},
		}, {
			Name: "defaultMaxDepth",
			Prompt: &survey.Select{
				Message: "Default max depth",
				Options: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				Default: strconv.Itoa(config.DefaultMaxDepth),
				Help:    "The default max depth for the roots",
			},
		}, {
			Name: "mostChosenCount",
			Prompt: &survey.Select{
				Message: "Most chosen count",
				Options: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				Default: strconv.Itoa(config.MostChosenCount),
				Help:    "The number of most chosen folders to display in the favorites list",
			},
		},
	}

	var err error
	if err = survey.Ask(questions, &values); err != nil {
		return nil, err
	}

	var (
		syncIntervalAnswer    = values["syncInterval"].(core.OptionAnswer)
		defaultMaxDepthAnswer = values["defaultMaxDepth"].(core.OptionAnswer)
		mostChosenCountAnswer = values["mostChosenCount"].(core.OptionAnswer)
	)

	config.SyncInterval, err = time.ParseDuration(syncIntervalAnswer.Value)
	if err != nil {
		return nil, err
	}

	config.DefaultMaxDepth, err = strconv.Atoi(defaultMaxDepthAnswer.Value)
	if err != nil {
		return nil, err
	}

	config.MostChosenCount, err = strconv.Atoi(mostChosenCountAnswer.Value)
	if err != nil {
		return nil, err
	}

	return config, nil
}
