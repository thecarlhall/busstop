package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Config is the configuration for the app
type Config struct {
	AppID     string
	Stops     []Stop
	Debug     bool
	Visual    bool
	Frequency int
}

// Stop is the stop information to report
type Stop struct {
	LocID     int
	Routes    []int
	Schedules []string
}

// LoadDefaultConfig loads the configuration file at ~/.busstop then overrides with CLI flags
func LoadConfig() (Config, error) {
	usr, _ := user.Current()
	configFilename := filepath.Join(usr.HomeDir, ".busstop")
	configFile, err := os.Open(configFilename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	dec := json.NewDecoder(configFile)
	dec.Decode(&config)

	if config.Debug {
		log.Printf("Config: %+v\n", config)
	}

	if config.validate() {
		return config, nil
	} else {
		return config, errors.New("Invalid configuration")
	}
}

func (config Config) validate() bool {
	valid := true

	if len(config.AppID) == 0 {
		log.Print("Config: appID is required\n")
		valid = false
	}

	if len(config.Stops) == 0 {
		log.Print("Config: stops are required\n")
		valid = false
	}

	return valid
}
