package main

import (
	"fmt"
	"log"
	"os/exec"
)

func displayMessage(title string, message string, showGrowl bool) {
	if showGrowl {
		script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
		cmd := exec.Command("/usr/bin/osascript", "-e", script)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(title)
		fmt.Println(message)
	}
}

func sprintRouteInfo(rs *ResultSet) string {
	var msg string
	for _, arrival := range rs.Arrival {
		msg += fmt.Sprintf("%-60s [%s]\n", arrival.FullSign, arrival.arrivalTime())
	}
	return msg
}

func main() {
	config := ParseFlags(nil)

	if *config.help {
		config.printHelp()
		return
	}

	rs := NewTrimetService(*config.appID, false).fetchLocationData(*config.locationID, *config.route)
	title := fmt.Sprintf("%60s\n", fmt.Sprintf("---[ Information for stop %d ]---", *config.locationID))
	routeInfo := sprintRouteInfo(rs)
	displayMessage(title, routeInfo, *config.growl)
}
