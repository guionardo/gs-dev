package calendar

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/configs"
)

func Subscribe(name string, uri string) error {
	if len(uri) == 0 {
		return fmt.Errorf("missing --uri argument")
	}
	if len(name) == 0 {
		return fmt.Errorf("missing --name argument")
	}

	cfg := getConfig()

	if cal := getCalendarByUri(uri); cal != nil {
		return fmt.Errorf("calendar just exists: %s", cal.Name)
	}

	if _, err := fetchUrl(uri); err != nil {
		return err
	}

	cal := &config.CalendarConfig{
		Name: name,
		URI:  uri,
	}
	if err := syncCacheCalendar(cfg, cal); err != nil {
		return err
	}
	cfg.Calendars[name] = cal
	return cfg.Save()

}

func Unsubscribe(calendarName string) error {
	cfg := getConfig()
	if _, ok := cfg.Calendars[calendarName]; !ok {
		return fmt.Errorf("calendar %s does not exists", calendarName)
	}
	if !configs.Confirm("Confirm unsubscribing of "+calendarName+" calendar", false) {
		return nil
	}
	delete(cfg.Calendars, calendarName)
	if err := cfg.Save(); err == nil {
		fmt.Printf("Calendar %s was unsubscribed\n", calendarName)
	} else {
		return fmt.Errorf("failed to unsubscribe calendar %s - %v", calendarName, err)
	}
	return nil
}

func syncCacheCalendar(cfg *config.CalendarsConfig, cal *config.CalendarConfig) error {
	force := false
	if data, err := cfg.GetCache(cal); len(data) == 0 || err != nil || !cfg.IsCalendarUpdated(cal) {
		force = true
	}
	if !force {
		return nil
	}
	data, err := fetchUrl(cal.URI)
	if err != nil {
		return err
	}
	if err = cfg.SetCache(cal, data); err == nil {
		cal.LastFetch = time.Now()
		err = cfg.Save()
	}

	return err
}

func fetchUrl(uri string) (data []byte, err error) {
	var resp *http.Response
	resp, err = http.DefaultClient.Get(uri)
	if resp != nil && resp.StatusCode >= 400 {
		err = fmt.Errorf("request failed to %s : %s", uri, resp.Status)
	}
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return
}

func fetch(cfg *config.CalendarsConfig, calendar *config.CalendarConfig) (cal *gocal.Gocal, err error) {
	if err = syncCacheCalendar(cfg, calendar); err != nil {
		return
	}
	data, err := cfg.GetCache(calendar)
	if err != nil {
		return
	}
	reader := bytes.NewReader(data)
	cal = gocal.NewParser(reader)
	start, end := time.Now().Add(-cfg.RangeInit), time.Now().Add(cfg.RangeEnd)
	cal.Start = &start
	cal.End = &end
	err = cal.Parse()
	return
}
