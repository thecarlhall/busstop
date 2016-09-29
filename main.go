package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// message creates and communicates (print, growl) a messages for the given parameters
func makeMessage(service *TrimetService, locID int, routes []int, schedules []string) Message {
	rs := service.FetchLocationData(locID, routes, schedules)

	title := fmt.Sprintf("Information For Stop %d", locID)
	var messages []string
	if len(rs.Arrival) > 0 {
		for _, arrival := range rs.Arrival {
			messages = append(messages, fmt.Sprintf("%-60s | %s %s\n", arrival.FullSign, arrival.ScheduledTime(), arrival.UntilArrival()))
		}
	}

	return Message{
		Subject: title,
		Bodies:  messages,
	}
}

func setupLogFile() *os.File {
	// Setup logging to go to a file
	usr, _ := user.Current()
	f, err := os.OpenFile(filepath.Join(usr.HomeDir, ".busstop.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(f)

	return f
}

func main() {
	logFile := setupLogFile()
	defer logFile.Close()

	config := LoadDefaultConfig()

	messenger := NewMessenger(config)
	for {
		service := NewTrimetService(config.AppID, config.Debug)
		for _, stop := range config.Stops {
			messages := makeMessage(service, stop.LocID, stop.Routes, stop.Schedules)
			messenger.Emit(messages)
		}

		if config.Frequency == 0 {
			break
		}

		time.Sleep(time.Duration(config.Frequency) * time.Second)
	}
}
