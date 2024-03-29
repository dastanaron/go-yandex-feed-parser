package main

import (
	"app/helpers"
	"app/types"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

const DEFAULT_CSV_COMMA = ';'

func main() {
	startTime := time.Now()

	appParameters := helpers.InitAppParams()

	file, err := os.Open(appParameters.XmlFilePath)
	if err != nil {
		helpers.CheckError("file not found", err)
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)

	localities := make(map[string]int)
	villages := make(map[string]int)
	categories := make(map[string]int)

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

				if offer.Category != "" {
					if _, ok := categories[offer.Category]; !ok {
						categories[offer.Category] = 0
					}
					categories[offer.Category]++
				}

				offersCount++
			}
		}
	}

	if appParameters.LocalitiesPathToSave != "" {
		saveData(appParameters.LocalitiesPathToSave, appParameters.FormatFile, localities)
	}

	if appParameters.VillagesPathToSave != "" {
		saveData(appParameters.VillagesPathToSave, appParameters.FormatFile, villages)
	}

	if appParameters.CategoriesPathToSave != "" {
		saveData(appParameters.CategoriesPathToSave, appParameters.FormatFile, categories)
	}

	endTime := time.Now()

	fmt.Println("Offers count: ", offersCount)
	fmt.Println("Unique localities count: ", len(localities))
	fmt.Println("Unique villages count: ", len(villages))
	fmt.Println("Unique object categories: ", len(categories))
	fmt.Println("All images count: ", imagesCount)

	fmt.Println("Time:", endTime.Sub(startTime))
}

func saveData(pathToSave, format string, data map[string]int) {
	var err error

	preparedPath := fmt.Sprintf("%s.%s", pathToSave, format)

	switch format {
	case "csv":
		err = helpers.WriteCsvData(data, preparedPath, DEFAULT_CSV_COMMA)
	case "xlsx":
		err = helpers.WriteXlsxData(data, preparedPath)
	default:
		fmt.Println("Format: ", format)
		fmt.Println("Undefined format for saving")
	}

	if err != nil {
		helpers.CheckError("Cannot write data", err)
	}
}
