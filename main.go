package main

import (
	"fmt"
	"log"
	"os/exec"
)

// displayMessage shows a growl message
func displayMessage(title string, message string) {
	script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
	cmd := exec.Command("/usr/bin/osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// printMessage prints a title and message to stdout
func printMessage(title string, message string) {
	fmt.Printf("%s\n%s", title, message)
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
func message(service *TrimetService, appID string, locID, route int, growl bool) {
	rs := service.FetchLocationData(locID, route)
	title := fmt.Sprintf("%60s", fmt.Sprintf("---[ Information For Stop %d ]---", locID))
	message := sprintRouteInfo(rs)
	if growl {
		displayMessage(title, message)
	} else {
		printMessage(title, message)
	}
}

func main() {
	config := LoadDefaultConfig()
	if config.Debug {
		fmt.Printf("%+v\n", config)
	}

	if config.Help {
		config.printHelp()
		return
	}

	service := NewTrimetService(config.AppID, config.Debug)
	for _, stop := range config.Stops {
		if len(stop.Routes) == 0 {
			message(service, config.AppID, stop.LocID, 0, config.Growl)
		} else {
			for _, route := range stop.Routes {
				message(service, config.AppID, stop.LocID, route, config.Growl)
			}
		}
	}
}
