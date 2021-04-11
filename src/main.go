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

	flag.StringVar(&xmlFilePath, "i", "../data/example.xml", "Path to XML feed")
	flag.StringVar(&localitiesPathToSave, "lo", "../data/result.csv", "Path for result about localities")

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

	helpers.WriteCsvLocalities(localities, localitiesPathToSave)

	endTime := time.Now()

	fmt.Println("Offers count:", offersCount)
	fmt.Println("Unique localities count", len(localities))

	fmt.Println("Time:", endTime.Sub(startTime))
}
