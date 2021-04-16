package main

import (
	"app/helpers"
	"app/types"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"time"
)

const DEFAULT_CSV_COMMA = ';'

func main() {
	startTime := time.Now()

	var xmlFilePath, localitiesPathToSave, villagesPathToSave, format string

	flag.StringVar(&xmlFilePath, "i", "", "[Required] Path to XML feed")
	flag.StringVar(&localitiesPathToSave, "lo", "", "[Optional] Path for result about localities. If empty, localities will not be saved")
	flag.StringVar(&villagesPathToSave, "vo", "", "[Optional] Path for result about villages. If empty, villages will not be saved")
	flag.StringVar(&format, "f", "csv", "[Optional] Format for saving file, default: csv")

	flag.Parse()

	file, err := os.Open(xmlFilePath)
	if err != nil {
		helpers.CheckError("file not found", err)
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)

	localities := make(map[string]int)
	villages := make(map[string]int)

	offersCount := 0
	imagesCount := 0

	for {
		tok, _ := decoder.Token()
		if tok == nil {
			break
		}

		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "offer" {
				var offer types.Offer
				decoder.DecodeElement(&offer, &tp)

				imagesCount += len(offer.Images)

				if offer.Location.LocalityName != "" {
					if _, ok := localities[offer.Location.LocalityName]; !ok {
						localities[offer.Location.LocalityName] = 0
					}
					localities[offer.Location.LocalityName]++
				}

				if offer.Location.VillageName != "" {
					if _, ok := villages[offer.Location.VillageName]; !ok {
						villages[offer.Location.VillageName] = 0
					}
					villages[offer.Location.VillageName]++
				}

				offersCount++
			}
		}
	}

	if localitiesPathToSave != "" {
		switch format {
		case "csv":
			err = helpers.WriteCsvData(localities, localitiesPathToSave, DEFAULT_CSV_COMMA)
		case "xlsx":
			err = helpers.WriteXlsxData(localities, localitiesPathToSave)
		default:
			fmt.Println("Undefined format for saving")
		}

		if err != nil {
			helpers.CheckError("Cannot write localities", err)
		}
	}

	if villagesPathToSave != "" {
		switch format {
		case "csv":
			err = helpers.WriteCsvData(villages, villagesPathToSave, DEFAULT_CSV_COMMA)
		case "xlsx":
			err = helpers.WriteXlsxData(villages, villagesPathToSave)
		default:
			fmt.Println("Undefined format for saving")
		}

		if err != nil {
			helpers.CheckError("Cannot write villages", err)
		}
	}

	endTime := time.Now()

	fmt.Println("Offers count:", offersCount)
	fmt.Println("Unique localities count", len(localities))
	fmt.Println("Unique villages count", len(villages))
	fmt.Println("All images count: ", imagesCount)

	fmt.Println("Time:", endTime.Sub(startTime))
}
