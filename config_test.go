package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	expected := Config{
		AppID: "abc123",
		Stops: []Stop{
			Stop{
				LocID:  123,
				Routes: []int{1, 2},
			},
			Stop{
				LocID:  124,
				Routes: []int{3},
			},
		},
		Debug: true,
		Growl: true,
		Help:  true,
	}

	tempFile, _ := ioutil.TempFile("", "busstop")
	defer os.Remove(tempFile.Name())

	enc := json.NewEncoder(tempFile)
	enc.Encode(expected)

	found := LoadConfig(tempFile.Name())

	if !reflect.DeepEqual(expected, *found) {
		t.Errorf("Expected %+v, found %+v", expected, found)
	}
}
