package main

import (
	"encoding/json"
	"flag"
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
	Help      bool
	Frequency int
}

// Stop is the stop information to report
type Stop struct {
	LocID     int
	Routes    []int
	Schedules []string
}

// LoadDefaultConfig loads the configuration file at ~/.busstop then overrides with CLI flags
func LoadDefaultConfig() *Config {
	usr, _ := user.Current()
	defaultConfigFile := filepath.Join(usr.HomeDir, ".busstop")

	config := LoadConfig(defaultConfigFile)
	if config.Debug {
		fmt.Printf("%+v\n", config)
	}

	return config
}

// LoadConfig loads the configuration file then overrides with CLI flags
func LoadConfig(defaultConfigFile string) *Config {
	configFile := flag.String("config", defaultConfigFile, "Config file to use")
	appID := flag.String("appID", "", "Trimet application ID")
	growl := flag.Bool("growl", false, "whether to use growl notifications")
	help := flag.Bool("help", false, "Show help information")
	debug := flag.Bool("debug", false, "Set debug mode")
	frequency := flag.Int("frequency", 0, "How often in minutes to poll when the requested data"+
		" is found.  If not set, the program runs once then exits.")

	flag.Parse()

	var config Config
	file, _ := ioutil.ReadFile(*configFile)
	json.Unmarshal(file, &config)

	if len(*appID) > 0 {
		config.AppID = *appID
	}

	if *debug {
		config.Debug = *debug
	}

	if *growl {
		config.Growl = *growl
	}

	if *help {
		config.Help = *help
	}

	if config.Debug {
		fmt.Printf("Loaded config file [%s]\n", *configFile)
	}

	if *frequency > 0 {
		config.Frequency = *frequency
	}

	config.validate()

	return &config
}

func (config *Config) validate() {
	if len(config.AppID) == 0 {
		log.Fatal("appID is required")
	}

	if len(config.Stops) == 0 {
		log.Fatal("stops are required")
	}
}

func (config *Config) PrintHelp() {
	flag.PrintDefaults()
}
