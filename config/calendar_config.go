package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	DefaultRangeInit = time.Duration(15) * time.Minute
	DefaultRangeEnd  = time.Hour * 24 * 7
)

//go:generate go run ../gen/configs.go
type CalendarsConfig struct {
	fileName      string
	lock          sync.Mutex
	Calendars     map[string]*CalendarConfig `yaml:"calendars"`
	FetchInterval time.Duration              `yaml:"fetch_interval"`
	RangeInit     time.Duration              `yaml:"range_init"`
	RangeEnd      time.Duration              `yaml:"range_end"`
	Enabled       bool                       `yaml:"enabled"`
}
type CalendarConfig struct {
	Name      string    `yaml:"name"`
	URI       string    `yaml:"uri"`
	LastFetch time.Time `yaml:"last_fetch"`
}

func NewCalendarsConfig(root string) (cfg *CalendarsConfig) {
	return &CalendarsConfig{
		fileName:      path.Join(root, "calendars.yaml"),
		Calendars:     make(map[string]*CalendarConfig, 0),
		FetchInterval: time.Hour,
		RangeInit:     DefaultRangeInit,
		RangeEnd:      DefaultRangeEnd,
	}
}

func (cfg *CalendarsConfig) IsCalendarUpdated(cal *CalendarConfig) bool {
	return time.Since(cal.LastFetch) <= cfg.FetchInterval
}
func (cfg *CalendarsConfig) CalendarCacheFile(cal *CalendarConfig) string {
	hash := md5.Sum([]byte(cal.URI))
	return path.Join(path.Dir(cfg.fileName), "cal_"+hex.EncodeToString(hash[:])+".txt")
}

func (cfg *CalendarsConfig) SetCache(cal *CalendarConfig, data []byte) error {
	cacheFile := cfg.CalendarCacheFile(cal)
	if len(data) == 0 {
		if _, err := os.Stat(cacheFile); err == nil {
			err = os.Remove(cacheFile)
			if err != nil {
				err = fmt.Errorf("failed removing calendar cache: %v", err)
			}
			return err
		}
	}
	file, err := os.Create(cacheFile)
	if err == nil {
		defer file.Close()
		_, err = file.Write(data)
	}
	if err != nil {
		err = fmt.Errorf("failed to set calendar cache: %v", err)
	}
	return err
}

func (cfg *CalendarsConfig) GetCache(cal *CalendarConfig) (data []byte, err error) {
	cacheFile := cfg.CalendarCacheFile(cal)
	var fileLength int64
	if stat, err := os.Stat(cacheFile); err != nil || stat.Size() == 0 {
		return data, fmt.Errorf("cache file not found or empty")
	} else {
		fileLength = stat.Size()
	}
	file, err := os.Open(cacheFile)
	if err != nil {
		err = fmt.Errorf("failed to open cache file: %v", err)
		return
	}
	data = make([]byte, fileLength)
	readen, err := file.Read(data)
	if readen != int(fileLength) {
		err = fmt.Errorf("failed to read cache file - expected %d bytes - got %d", fileLength, readen)
	}
	return
}
