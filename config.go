package main

import (
	"flag"
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
func ParseFlags() *Config {
	growl := flag.Bool("growl", false, "whether to use growl notifications")
	appID := flag.String("appID", "", "Trimet application ID")
	locationID := flag.Int("locationID", 13168, "location ID to track")
	route := flag.Int("route", 0, "Route number to filter by")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		printHelp()
		return &Config{}
	}

	if len(*appID) == 0 {
		log.Fatal("appID is required")
	}

	return &Config{
		appID:      appID,
		locationID: locationID,
		route:      route,
		growl:      growl,
		help:       help,
	}
}
