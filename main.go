package main

import (
	"flag"
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

func printRouteInfo(locationID int, rs *ResultSet) {
	fmt.Printf("%60s\n", fmt.Sprintf("---[ Information for stop %d ]---", locationID))
	for _, arrival := range rs.Arrival {
		fmt.Printf("%-60s [%s]\n", arrival.FullSign, arrival.arrivalTime())
	}
}

func printHelp() {
	fmt.Println("  busstop [OPTIONS]")
	fmt.Println("")
	fmt.Println("  Required")
	fmt.Println("    --appID <app_id>")
	fmt.Println("    --locationID <loc_id>")
	fmt.Println("")
	fmt.Println("  Optional")
	fmt.Println("    --route <route>")
	fmt.Println("    --help")
}

func main() {
	//growl := flag.Bool("growl", false, "whether to use growl notifications")
	appID := flag.String("appID", "", "Trimet application ID")
	locationID := flag.Int("locationID", 13168, "location ID to track")
	route := flag.Int("route", 0, "Route number to filter by")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if len(*appID) == 0 {
		log.Fatal("appID is required")
	}

	rs := NewTrimetService(*appID, false).fetchLocationData(*locationID, *route)
	printRouteInfo(*locationID, rs)
}
