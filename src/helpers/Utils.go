package helpers

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"

	xlsx "github.com/tealeg/xlsx/v3"
)

func WriteXlsxData(data map[string]int, pathToSave string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	preparedData := prepareData(data)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet_1")
	if err != nil {
		return err
	}

	sheet.SetColWidth(1, 1, 35)

	headStyleCell := xlsx.Style{
		Font: xlsx.Font{
			Bold: true,
			Name: "Verdana",
		},
		Alignment: xlsx.Alignment{
			Horizontal: "center",
		},
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "name"
	cell.SetStyle(&headStyleCell)

	cell = row.AddCell()
	cell.Value = "count"
	cell.SetStyle(&headStyleCell)

	for _, cells := range preparedData {
		row = sheet.AddRow()

		for _, cellValue := range cells {
			cell = row.AddCell()
			cell.Value = cellValue
		}
	}

	err = file.Save(pathToSave)
	if err != nil {
		return err
	}

	return nil
}

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

func prepareData(namesMap map[string]int) [][]string {
	outputData := make([][]string, 0)

	keys := make([]string, 0, len(namesMap))
	for k := range namesMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		data := []string{k, strconv.Itoa(namesMap[k])}
		outputData = append(outputData, data)
	}

	return outputData
}

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
