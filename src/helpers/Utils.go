package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
)

func WriteCsvData(data map[string]int, pathToSave string, comma rune) error {
	resultFile, err := os.Create(pathToSave)
	CheckError("Cannot create file", err)

	writer := csv.NewWriter(resultFile)

	writer.Comma = comma
	err = writer.WriteAll(prepareData(data))

	if err != nil {
		return err
	}

	writer.Flush()

	err = resultFile.Close()

	if err != nil {
		return err
	}

	return nil
}

func prepareData(localities map[string]int) [][]string {
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
