package parser

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"web-scraper/handlers"
)

type ChanUrls struct {
	url   string
	bytes int
	time  string
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

		// fmt.Println("HTML:\n\n", string(bytes))
		c <- ChanUrls{
			bytes: len(out),
			url:   record,
			time:  track,
		}
		// fmt.Println(<-c)
	}
}
