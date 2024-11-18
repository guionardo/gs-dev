package calendar

import "github.com/guionardo/gs-dev/config"

var readenConfig *config.CalendarsConfig

func getConfig() *config.CalendarsConfig {
	if readenConfig == nil {
		readenConfig = config.NewCalendarsConfig(config.GetConfigRepositoryFolder())
		_ = readenConfig.Load()
	}

	return readenConfig
}

func getCalendarByUri(uri string) *config.CalendarConfig {
	for _, cal := range getConfig().Calendars {
		if cal.URI == uri {
			return cal
		}
	}
	return nil
}
