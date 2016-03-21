package main

import (
	"flag"
	"fmt"
	"log"
)

// Config is the configuration for the app
type Config struct {
	growl      *bool
	appID      *string
	locationID *int
	route      *int
	help       *bool
}

// ParseFlags parses cli flags to create a config
func ParseFlags(config *Config) *Config {
	growl := flag.Bool("growl", false, "whether to use growl notifications")
	appID := flag.String("appID", "", "Trimet application ID")
	locationID := flag.Int("locationID", 0, "location ID to track")
	route := flag.Int("route", 0, "Route number to filter by")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if config == nil {
		return &Config{
			appID:      appID,
			locationID: locationID,
			route:      route,
			growl:      growl,
			help:       help,
		}
	}

	if len(*appID) > 0 {
		config.appID = appID
	}

	if *locationID > 0 {
		config.locationID = locationID
	}

	if *route > 0 {
		config.route = route
	}

	if *growl {
		config.growl = growl
	}

	if *help {
		config.help = help
	}

	return config
}

func (config *Config) validate() {
	if len(*config.appID) == 0 {
		log.Fatal("appID is required")
	}

	if *config.locationID == 0 {
		log.Fatal("locationID is required")
	}
}

func (config *Config) printHelp() {
	fmt.Println("  busstop [OPTIONS]")
	fmt.Println("")
	fmt.Println("  Required")
	fmt.Println("    --appID <app_id>")
	fmt.Println("    --locationID <loc_id>")
	fmt.Println("")
	fmt.Println("  Optional")
	fmt.Println("    --route <route>")
	fmt.Println("    --growl")
	fmt.Println("    --help")
}
