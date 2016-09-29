package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
)

// Config is the configuration for the app
type Config struct {
	AppID     string
	Stops     []Stop
	Debug     bool
	Growl     bool
	Frequency int
}

// Stop is the stop information to report
type Stop struct {
	LocID     int
	Routes    []int
	Schedules []string
}

// LoadDefaultConfig loads the configuration file at ~/.busstop then overrides with CLI flags
func LoadDefaultConfig() Config {
	usr, _ := user.Current()
	configFile := filepath.Join(usr.HomeDir, ".busstop")

	var config Config
	file, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(file, &config)

	if config.Debug {
		fmt.Printf("%+v\n", config)
	}

	config.validate()

	return config
}

func (config Config) validate() {
	if len(config.AppID) == 0 {
		log.Fatal("appID is required")
	}

	if len(config.Stops) == 0 {
		log.Fatal("stops are required")
	}
}
