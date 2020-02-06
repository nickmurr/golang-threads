package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
	"web-scraper/handlers"
	"web-scraper/parser"
)

var (
	maxGoroutines = 100
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
	//

	out := WriteToFile(len(records), c)
	fmt.Println(out)

}

func WriteToFile(length int, c chan parser.ChanUrls) []parser.ChanUrls {
	// var out []parser.ChanUrls
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	// var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	for i := 0; i < length; i++ {
		x:=<-c
		row = sheet.AddRow()

		_URL := row.AddCell()
		_Size := row.AddCell()
		_Time := row.AddCell()

		_URL.Value = x.Url
		_Size.Value = fmt.Sprint(x.Bytes)
		_Time.Value = fmt.Sprint(x.Time)
	}

	err = file.Save(fmt.Sprintf("./output/MyXLSXFile-%v.xlsx", time.Now().Unix()))
	if err != nil {
		fmt.Printf(err.Error())
	}

	return nil

}
