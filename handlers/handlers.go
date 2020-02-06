package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

func GetRecords() []string {
	xlFile, err := xlsx.OpenFile("input2.xlsx")
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

func TimeTrack(start time.Time, print bool, message string) string {
	elapsed := time.Since(start)
	if len(message) == 0 {
		message = "Total"
	}
	if print {
		fmt.Printf("\n%s took %s", message, elapsed)
	}
	return fmt.Sprint(elapsed)
}
