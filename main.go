package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	growl := flag.Bool("growl", false, "whether to use growl notifications")
	locationId := flag.Int("locationId", 13168, "location ID to track")
	flag.Parse()

	doc, err := goquery.NewDocument(fmt.Sprintf("http://trimet.org/arrivals/small/tracker?locationID=%d", *locationId))
	if err != nil {
		log.Fatal(err)
	}

	title := fmt.Sprintf("Bus Stop - %d", *locationId)

	selection := doc.Find("ul#arrivalslist.group > li")
	messages := make([]string, selection.Length())
	selection.Each(func(i int, s *goquery.Selection) {
		messages[i] = fmt.Sprint(strings.TrimSpace(s.Find("h3").Text()), "\n    ", s.Find("p.clear").Text(), " in ", s.Find("p.arrival").Text())
	})

	message := strings.Join(messages, "\n")
	if *growl {
		script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
		cmd := exec.Command("/usr/bin/osascript", "-e", script)
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(title)
		fmt.Println(message)
	}
}
