package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"web-scraper/handlers"
	"web-scraper/parser"
)

var (
	maxGoroutines = 100
)

func main() {
	newpath := filepath.Join(".", "dist")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	guard := make(chan struct{}, maxGoroutines)
	defer handlers.TimeTrack(time.Now(), true, "Total scraping")

	records := handlers.GetRecords()
	c := make(chan parser.ChanUrls, len(records))

	for _, v := range records {
		guard <- struct{}{}
		go func(v string, c chan parser.ChanUrls) {
			parser.ReadBody(v, c)
			<-guard
		}(v, c)

	}

	parser.WriteToFile(len(records), c)

}
