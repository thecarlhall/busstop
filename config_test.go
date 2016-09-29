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
			{
				LocID:  123,
				Routes: []int{1, 2},
			},
			{
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

	json.NewEncoder(tempFile).Encode(expected)
	tempFile.Close()

	found := LoadConfig(tempFile.Name())

	if !reflect.DeepEqual(expected, *found) {
		t.Errorf("Expected %+v, found %+v", expected, found)
	}
}
