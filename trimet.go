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

func (arrival *Arrival) ScheduledTime() string {
	scheduled := time.Unix(0, int64(arrival.Scheduled)*int64(time.Millisecond))
	return scheduled.Format("3:04pm")
}

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

// TrimetService handles calls to Trimet's API
type TrimetService struct {
	baseURL string
	appID   string
	debug   bool
}

// NewTrimetService is the creator for TrimetService
func NewTrimetService(appID string, debug bool) *TrimetService {
	return &TrimetService{
		baseURL: "https://developer.trimet.org/ws/v2/arrivals",
		appID:   appID,
		debug:   debug,
	}
}

func (ts *TrimetService) fetchLocationData(locID int, route int) *ResultSet {
	// output for remote call
	rs := &ArrivalsResponse{}
	url := fmt.Sprintf("%s?appID=%s&locIDs=%d", ts.baseURL, ts.appID, locID)
	ts.getJSON(url, rs)

	var arrivals []Arrival
	if route > 0 {
		for _, arrival := range rs.ResultSet.Arrival {
			if arrival.Route == route {
				arrivals = append(arrivals, arrival)
			}
		}
		rs.ResultSet.Arrival = arrivals
	}
	return &rs.ResultSet
}

func (ts *TrimetService) getJSON(url string, target interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatal(fmt.Sprintf("Unable to get data: %s", resp.Status))
	}

	if ts.debug {
		body := resp.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		s := buf.String()
		fmt.Printf("Body: %s\n", s)
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Fatal(err)
	}
}
