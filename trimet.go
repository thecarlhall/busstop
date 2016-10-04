package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ArrivalsResponse holds the result set from the arrivals api
type ArrivalsResponse struct {
	ResultSet ResultSet
}

// ResultSet is the container of a response
type ResultSet struct {
	Arrival   []Arrival
	Detour    []Detour
	Location  []Location
	QueryTime int
}

// Arrival is the model of an arrival at a stop.
type Arrival struct {
	BlockID      int
	Departed     bool
	Detour       []string
	Detoured     bool
	Dir          int
	Estimated    int
	Feet         int
	FullSign     string
	ID           string
	InCongestion bool
	LocID        int
	NewTrip      bool
	Piece        string
	Route        int
	Scheduled    int
	ShortSign    string
	Status       string
	TripID       string
	VehicleID    string
}

// ScheduledTime returns a string with the time of the next arrival
func (arrival *Arrival) ScheduledTime() string {
	scheduled := time.Unix(0, int64(arrival.Scheduled)*int64(time.Millisecond))
	return scheduled.Format("3:04pm")
}

// UntilArrival returns a string with the time until the next arrival
func (arrival *Arrival) UntilArrival() string {
	scheduled := time.Unix(0, int64(arrival.Scheduled)*int64(time.Millisecond))
	estimated := time.Unix(0, int64(arrival.Estimated)*int64(time.Millisecond))
	duration := -time.Since(estimated).Minutes()

	var durationStr string

	if duration > 60 || arrival.Estimated == 0 {
		durationStr = "on " + scheduled.Format("Mon, 02 Jan 2006")
	} else if duration < 1 {
		durationStr = "Due now!"
	} else {
		durationStr = "in " + strconv.FormatFloat(duration, 'f', 0, 64) + "m"
	}

	return durationStr
}

// Route is the data for a transit route
type Route struct {
	Desc   string
	Detour bool
	Route  int
	Type   string
}

// Detour is the data for a route detour
type Detour struct {
	Begin       int
	Desc        string
	End         int
	ID          string
	InfoLinkURL string `json:"info_link_url"`
	Route       []Route
}

// Location is the data for a location on a route
type Location struct {
	Desc          string
	Dir           string
	ID            int
	Lat           float32
	Lng           float32
	PassengerCode string
}

// NewTrimetService is the creator for TrimetService
func NewTrimetService(config Config) *TrimetService {
	return &TrimetService{
		BaseURL: "https://developer.trimet.org/ws/v2/arrivals",
		AppID:   config.AppID,
		Debug:   config.Debug,
	}
}

// TrimetService handles calls to Trimet's API
type TrimetService struct {
	BaseURL string
	AppID   string
	Debug   bool
}

// FetchLocationData fetches location data from the Trimet API
func (ts *TrimetService) FetchLocationData(locID int, routes []int, schedules []string) *ResultSet {
	// output for remote call
	rs := &ArrivalsResponse{}
	url := fmt.Sprintf("%s?appID=%s&locIDs=%d", ts.BaseURL, ts.AppID, locID)
	ts.getJSON(url, rs)

	if len(routes) > 0 {
		var arrivals []Arrival
		for _, arrival := range rs.ResultSet.Arrival {
			for _, route := range routes {
				if arrival.Route == route {
					arrivals = append(arrivals, arrival)
				}
			}
		}
		rs.ResultSet.Arrival = arrivals
	}

	if len(schedules) > 0 {
		var arrivals []Arrival
		for _, arrival := range rs.ResultSet.Arrival {
			for _, schedule := range schedules {
				if arrival.ScheduledTime() == schedule {
					arrivals = append(arrivals, arrival)
				}
			}
		}
		rs.ResultSet.Arrival = arrivals
	}

	return &rs.ResultSet
}

func (ts *TrimetService) getJSON(url string, target interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Printf("Unable to get data: %s", resp.Status)
	}

	if ts.Debug {
		body := resp.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		s := buf.String()
		fmt.Printf("Body: %s\n", s)
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Print(err)
	}
}
