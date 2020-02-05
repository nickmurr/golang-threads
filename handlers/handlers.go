package handlers

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"time"
)

func GetRecords() []string {
	xlFile, err := xlsx.OpenFile("input1.xlsx")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	var records []string

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				if len(text) > 1 {
					records = append(records, strings.TrimSpace(text))
				}
			}
		}
	}

	return records
}

func TimeTrack(start time.Time, print bool) string {
	elapsed := time.Since(start)
	if print {
		fmt.Printf("\n%s took %s", "Total", elapsed)
	}
	return fmt.Sprint(elapsed)
}
