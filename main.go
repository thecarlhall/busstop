package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	locationId := 13168
	doc, err := goquery.NewDocument(fmt.Sprintf("http://trimet.org/arrivals/small/tracker?locationID=%d", locationId))
	if err != nil {
		log.Fatal(err)
	}

	title := fmt.Sprintf("Bus Stop - %d", locationId)

	selection := doc.Find("ul#arrivalslist.group > li")
	messages := make([]string, selection.Length())
	selection.Each(func(i int, s *goquery.Selection) {
		messages[i] = fmt.Sprint(s.Find("p.clear").Text(), " in ", s.Find("p.arrival").Text())
	})

	message := strings.Join(messages, "\n")
	script := fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title)
	cmd := exec.Command("/usr/bin/osascript", "-e", script)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
