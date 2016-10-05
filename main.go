package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	arrivingBusTime = regexp.MustCompile(`in (\d+)m`)
)

// message creates and communicates (print, growl) a messages for the given parameters
func makeMessage(service *TrimetService, stop Stop, notitificationThreshold int) Message {
	rs := service.FetchLocationData(stop)

	title := fmt.Sprintf("Information For Stop %d", stop.LocID)
	var messages []string
	if len(rs.Arrival) > 0 {
		for _, arrival := range rs.Arrival {
			approxArrivalTime := arrival.UntilArrival()

			arrivalTime := arrivingBusTime.FindStringSubmatch(approxArrivalTime)
			if arrivalTime == nil {
				continue
			}

			minutesUntil, err := strconv.ParseInt(arrivalTime[1], 10, 32)
			if err != nil {
				continue
			}

			if int(minutesUntil) <= notitificationThreshold {
				messages = append(messages, fmt.Sprintf("%-60s | %s %s\n", arrival.FullSign, arrival.ScheduledTime(), approxArrivalTime))
			}
		}
	}

	return Message{
		Subject: title,
		Bodies:  messages,
	}
}

func setupLogFile() (*os.File, error) {
	// Setup logging to go to a file
	usr, _ := user.Current()
	f, err := os.OpenFile(filepath.Join(usr.HomeDir, ".busstop.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(f)

	return f, err
}

func main() {
	logFile, err := setupLogFile()
	if err != nil {
		log.Print(err)
		return
	}
	defer logFile.Close()

	config, err := LoadConfig()
	if err != nil {
		return
	}

	messenger := NewMessenger(config)
	for {
		service := NewTrimetService(config)
		for _, stop := range config.Stops {
			messages := makeMessage(service, stop, config.NotificationThreshold)
			messenger.Emit(messages)
		}

		if config.PollingFrequency <= 0 {
			break
		}

		time.Sleep(time.Duration(config.PollingFrequency) * time.Minute)
	}
}
