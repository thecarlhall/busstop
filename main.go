package main

import (
	"fmt"
	"log"
	"os/exec"
)

// displayNotification shows a growl message
func notify(title string, message string) {
	script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
	cmd := exec.Command("/usr/bin/osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// console prints a title and message to stdout
func console(title string, message string) {
	fmt.Printf("%60s\n%s\n", fmt.Sprintf("---[  %s  ]---", title), message)
}

// sprintRouteInfo creates a string that represents the stops of the result set
func sprintRouteInfo(rs *ResultSet) string {
	var msg string
	for _, arrival := range rs.Arrival {
		msg += fmt.Sprintf("%-60s | %s %s\n", arrival.FullSign, arrival.ScheduledTime(), arrival.UntilArrival())
	}
	return msg
}

// message creates and communicates (print, growl) a messages for the given parameters
func message(service *TrimetService, appID string, locID int, routes []int, schedules []string, growl bool) {
	rs := service.FetchLocationData(locID, routes, schedules)

	if len(rs.Arrival) > 0 {
		title := fmt.Sprintf("Information For Stop %d", locID)
		message := sprintRouteInfo(rs)
		if growl {
			notify(title, message)
		} else {
			console(title, message)
		}
	}
}

func main() {
	config := LoadDefaultConfig()

	if config.Help {
		config.printHelp()
		return
	}

	service := NewTrimetService(config.AppID, config.Debug)
	for _, stop := range config.Stops {
		message(service, config.AppID, stop.LocID, stop.Routes, stop.Schedules, config.Growl)
	}
}
