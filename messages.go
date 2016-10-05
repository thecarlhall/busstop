package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func NewMessenger(config Config) Messenger {
	if config.Display == "growl" {
		return GrowlMessenger{}
	} else {
		return ConsoleMessenger{}
	}
}

type Message struct {
	Subject string
	Bodies  []string
}

type Messenger interface {
	Emit(message Message)
}

// ConsoleMessager prints a formatted title and message to stdout for all message bodies
type ConsoleMessenger struct{}

func (ConsoleMessenger) Emit(message Message) {
	msg := strings.Join(message.Bodies, "\n")
	fmt.Printf("%60s\n%s\n", fmt.Sprintf("---[  %s  ]---", message.Subject), msg)
}

// GrowlMessages shows growl messages for each message body
type GrowlMessenger struct{}

func (GrowlMessenger) Emit(message Message) {
	// iterate in reverse so that the latest time is shown last and longest
	for i := len(message.Bodies) - 1; i >= 0; i-- {
		body := message.Bodies[i]
		script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", body, message.Subject)
		cmd := exec.Command("/usr/bin/osascript", "-e", script)
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}
	}
}
