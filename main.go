package main

import (
	"fmt"
	"time"
	"web-scraper/handlers"
	"web-scraper/parser"
)

var (
	maxGoroutines = 80
)

func main() {
	guard := make(chan struct{}, maxGoroutines)

	defer handlers.TimeTrack(time.Now(), true)

	records := handlers.GetRecords()
	c := make(chan parser.ChanUrls, len(records))

	for _, v := range records {
		guard <- struct{}{}
		go func(v string, c chan parser.ChanUrls) {
			parser.ReadBody(v, c)
			<-guard
		}(v, c)

	}

	for _ = range records {
		fmt.Println(<-c)
	}

}
