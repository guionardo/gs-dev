package calendar

import (
	"errors"
	"time"
)

func Setup(rangeInit time.Duration, rangeEnd time.Duration) error {
	cfg := getConfig()
	cfg.RangeInit = rangeInit
	cfg.RangeEnd = rangeEnd

	for _, cal := range cfg.Calendars {
		cal.LastFetch = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return cfg.Save()
}

func Enable() error {
	cfg := getConfig()
	if cfg.Enabled {
		return errors.New("calendar is just enabled")
	}
	cfg.Enabled = true
	return cfg.Save()
}

func Disable() error {
	cfg := getConfig()
	if !cfg.Enabled {
		return errors.New("calendar is just disabled")
	}
	cfg.Enabled = false
	return cfg.Save()
}
