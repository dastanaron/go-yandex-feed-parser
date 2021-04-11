package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
)

func WriteCsvLocalities(localities map[string]int, pathToSave string) {
	resultFile, err := os.Create(pathToSave)
	CheckError("Cannot create file", err)

	writer := csv.NewWriter(resultFile)

	for _, localityData := range prepareLocalities(localities) {
		err := writer.Write(localityData)
		CheckError("Cannot write to file", err)
	}
	writer.Flush()
	resultFile.Close()
}

func prepareLocalities(localities map[string]int) [][]string {
	localitiesData := make([][]string, 0)

	keys := make([]string, 0, len(localities))
	for k := range localities {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		data := []string{k, strconv.Itoa(localities[k])}
		localitiesData = append(localitiesData, data)
	}

	return localitiesData
}

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
