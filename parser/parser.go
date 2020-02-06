package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"web-scraper/handlers"

	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
)

type ChanUrls struct {
	Url   string
	Bytes string
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
		if err == io.ErrUnexpectedEOF {
			fmt.Println("Unexpected EOF")
			return
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()
		out := string(bytes)
		track := handlers.TimeTrack(t, false, "")

		fmt.Println(track)

		htmlBody := []byte(out)

		parse, err := url.Parse(record)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("dist/%s-%v.html", parse.Host, uuid.New()), htmlBody, 0666)
		if err != nil {
			log.Fatal(err)
		}
		c <- ChanUrls{
			Bytes: fmt.Sprint(len(out)),
			Url:   record,
			Time:  track,
		}

	}
}

func WriteToFile(length int, c chan ChanUrls) []ChanUrls {
	defer handlers.TimeTrack(time.Now(), true, "Writing to file")
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
		fmt.Println("Created Row")
		row = sheet.AddRow()

		_URL := row.AddCell()
		_Size := row.AddCell()
		_Time := row.AddCell()

		_URL.Value = x.Url
		_Size.Value = fmt.Sprint(x.Bytes)
		_Time.Value = fmt.Sprint(x.Time)
	}

	fmt.Println("Save to file")
	err = file.Save(fmt.Sprintf("./dist/MyXLSXFile-%v.xlsx", time.Now().Unix()))
	if err != nil {
		fmt.Println(<-c)
		fmt.Printf(err.Error())
	}

	return nil
}
