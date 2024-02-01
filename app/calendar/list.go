package calendar

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/apognu/gocal"
	"github.com/guionardo/gs-dev/config"
	"github.com/mitchellh/colorstring"
)

type EventItem struct {
	CalendarName string
	Event        gocal.Event
}

func fetchAll(cfg *config.CalendarsConfig) {
	wg := sync.WaitGroup{}
	wg.Add(len(cfg.Calendars))
	for _, cal := range cfg.Calendars {
		go func(c *config.CalendarConfig) {
			defer wg.Done()
			if _, err := fetch(cfg, c); err != nil {
				colorstring.Printf("[_red_]ERROR[reset] %v", err)
			}
		}(cal)
	}
	wg.Wait()
}

func List(noUriLink bool) error {
	cfg := getConfig()
	if !cfg.Enabled {
		return nil
	}
	fetchAll(cfg)
	eventList := make([]EventItem, 0, 10)
	maxNameLength := 0
	colors := []string{"[cyan]", "[magenta]", "[red]", "[light_blue]"}
	calendarColors := make(map[string]string)
	calIndex := 0
	for name, calendar := range cfg.Calendars {
		calendarColors[name] = colors[calIndex%len(colors)]
		calIndex++
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
		cal, err := fetch(cfg, calendar)
		if err != nil {
			continue
		}
		for _, event := range cal.Events {
			eventList = append(eventList, EventItem{
				CalendarName: name,
				Event:        event,
			})
		}
	}
	if len(eventList) == 0 {
		fmt.Println("No events on calendar")
		return nil
	}
	sort.Slice(eventList, func(i int, j int) bool {
		return eventList[i].Event.Start.Before(*eventList[j].Event.Start)
	})
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	lastEventDay := 0
	showNextEvent := false
	for _, event := range eventList {
		eventDateTime := event.Event.Start.In(loc)
		if eventDateTime.Day() == lastEventDay {
			fmt.Print("           ")
		} else {
			colorstring.Printf("[yellow]%s[default] ", eventDateTime.Format("02/01/2006"))
			lastEventDay = eventDateTime.Day()
		}

		timePrefix := "[green]"
		if !showNextEvent && eventDateTime.After(time.Now()) {
			timePrefix = "[_green_][black]"
			showNextEvent = true
		}
		colorstring.Printf(timePrefix+"%s", eventDateTime.Format("15:04"))

		calName := fmt.Sprintf(fmt.Sprintf(" %%-%ds", maxNameLength), event.CalendarName)
		colorstring.Printf(calendarColors[calName]+"[bold]%s[default]: ", calName)

		// fmt.Printf("%s on %s by %s %s\n", e.Summary, e.Start, e.Organizer.Cn, url)
		fmt.Printf("- %s\n", event.Event.Summary)
		if url, caption := extractLink(event.Event); len(url) > 0 {
			if !noUriLink {
				url = consoleUrl(url, caption)
			}
			fmt.Printf("%s %s\n", strings.Repeat(" ", maxNameLength+7), url)
		}

	}
	return nil
}

func extractLink(event gocal.Event) (url string, caption string) {
	if googleMeet, ok := event.CustomAttributes["X-GOOGLE-CONFERENCE"]; ok {
		url = googleMeet
		caption = "Google Meet"
		return
	}
	if p := strings.Index(event.Description, "<https://teams.microsoft.com/l/meetup-join"); p >= 0 {
		url = event.Description[p:]
		if p = strings.Index(url, ">"); p >= 0 {
			url, _ = strings.CutPrefix(url[0:p], "<")
			caption = "Teams Meet"
			return
		}
		url = ""
	}
	return
}
