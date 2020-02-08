package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"web-scraper/handlers"
	"web-scraper/parser"
)

var (
	maxGoroutines = 40
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	newpath := filepath.Join(".", "dist/html")
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
