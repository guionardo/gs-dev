package cmd

import (
	"github.com/guionardo/gs-dev/app/calendar"
	"github.com/guionardo/gs-dev/config"
	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Remote calendars",
	Long:  "List calendars itens, and manage calendars subscriptions",
}

func init() {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Settings for calendar",
		RunE:  setupCalendar,
	}
	setupCmd.Flags().DurationP("range_init", "i", config.DefaultRangeInit, "List events after (now - range_init)")
	setupCmd.Flags().DurationP("range_end", "e", config.DefaultRangeEnd, "List events before (now + range_end)")

	enableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable calendar listing",
		RunE:  enableCalendar,
	}

	disableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable calendar listing",
		RunE:  disableCalendar,
	}

	subscribeCalendarCmd := &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe remote calendar [URI]",
		RunE:  subscribeCalendar,
	}
	subscribeCalendarCmd.Flags().StringP("uri", "u", "", "Remote calendar URI (ics)")
	subscribeCalendarCmd.Flags().StringP("name", "n", "", "Calendar name")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all events",
		RunE:  calendarList,
	}
	listCmd.Flags().BoolP("no_uri_link", "n", false, "Disable console link for event URI")

	calendarCmd.AddCommand(
		setupCmd,
		enableCmd,
		disableCmd,
		subscribeCalendarCmd,
		listCmd,
	)

	rootCmd.AddCommand(calendarCmd)
}

func setupCalendar(cmd *cobra.Command, args []string) error {
	rangeInit, _ := cmd.Flags().GetDuration("range_init")
	rangeEnd, _ := cmd.Flags().GetDuration("range_end")

	return calendar.Setup(rangeInit, rangeEnd)
}

func calendarList(cmd *cobra.Command, args []string) error {
	noUriLink, _ := cmd.Flags().GetBool("no_uri_link")
	return calendar.List(noUriLink)
}

func subscribeCalendar(cmd *cobra.Command, args []string) error {
	uri, _ := cmd.Flags().GetString("uri")
	name, _ := cmd.Flags().GetString("name")
	return calendar.Subscribe(name, uri)
}

func enableCalendar(cmd *cobra.Command, args []string) error {
	return calendar.Enable()
}

func disableCalendar(cmd *cobra.Command, args []string) error {
	return calendar.Disable()
}
