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

func main() {
	startTime := time.Now()

	var xmlFilePath, localitiesPathToSave string

	flag.StringVar(&xmlFilePath, "i", "", "[Required] Path to XML feed")
	flag.StringVar(&localitiesPathToSave, "lo", "", "[Optional] Path for result about localities. If empty, localities will not be saved")

	flag.Parse()

	file, err := os.Open(xmlFilePath)
	if err != nil {
		helpers.CheckError("file not found", err)
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)

	localities := make(map[string]int)

	offersCount := 0
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

				if offer.Location.LocalityName != "" {
					if _, ok := localities[offer.Location.LocalityName]; !ok {
						localities[offer.Location.LocalityName] = 0
					}
					localities[offer.Location.LocalityName]++
				}

				offersCount++
			}
		}
	}

	if localitiesPathToSave != "" {
		helpers.WriteCsvLocalities(localities, localitiesPathToSave)
	}

	endTime := time.Now()

	fmt.Println("Offers count:", offersCount)
	fmt.Println("Unique localities count", len(localities))

	fmt.Println("Time:", endTime.Sub(startTime))
}
