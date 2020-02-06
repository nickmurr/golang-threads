package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"web-scraper/handlers"
)

type ChanUrls struct {
	Url   string
	Bytes int
	Time  string
}

func ReadBody(record string, c chan ChanUrls) {
	t := time.Now()
	if len(record) > 1 {
		resp, err := http.Get(record)
		if err != nil {
			log.Fatal(err)
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()
		out := string(bytes)
		track := handlers.TimeTrack(t, false)
		fmt.Println(track)
		c <- ChanUrls{
			Bytes: len(out),
			Url:   record,
			Time:  track,
		}
	}
}
