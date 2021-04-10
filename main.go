package main

import (
	"app/helpers"
	"app/types"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	file, err := os.Open("./data/example.xml")
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

	helpers.WriteCsvLocalities(localities)

	endTime := time.Now()

	fmt.Println("Offers count:", offersCount)
	fmt.Println("Unique localities count", len(localities))

	fmt.Println("Time:", endTime.Sub(startTime))
}
