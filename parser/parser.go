package parser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
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
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Error get url", err)
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

func WriteToFile(length int, c chan ChanUrls) []ChanUrls {
	defer handlers.TimeTrack(time.Now(),true)
	fmt.Println("started")
	// var out []parser.ChanUrls
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	// var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("URL Output")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()

	urlTitle := row.AddCell()
	sizeTitle := row.AddCell()
	timeTitle := row.AddCell()

	urlTitle.Value = "URL"
	sizeTitle.Value = "Size"
	timeTitle.Value = "Time"

	for i := 0; i < length; i++ {
		x := <-c
		row = sheet.AddRow()

		_URL := row.AddCell()
		_Size := row.AddCell()
		_Time := row.AddCell()

		_URL.Value = x.Url
		_Size.Value = fmt.Sprint(x.Bytes)
		_Time.Value = fmt.Sprint(x.Time)
	}

	err = file.Save(fmt.Sprintf("./dist/MyXLSXFile-%v.xlsx", time.Now().Unix()))
	if err != nil {
		fmt.Println(<-c)
		fmt.Printf(err.Error())
	}

	return nil
}
