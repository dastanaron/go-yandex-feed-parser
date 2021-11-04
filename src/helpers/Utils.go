package helpers

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	xlsx "github.com/tealeg/xlsx/v3"
)

type AppParameters struct {
	XmlFilePath          string
	LocalitiesPathToSave string
	VillagesPathToSave   string
	CategoriesPathToSave string
	FormatFile           string
	IsInteractive        bool
}

func InitAppParams() AppParameters {
	var xmlFilePath, localitiesPathToSave, villagesPathToSave, categoriesPathToSave, format string
	var isInteractive bool

	flag.StringVar(&xmlFilePath, "s", "", "[Required] Path to XML feed")
	flag.BoolVar(&isInteractive, "i", false, "[Optional] Interactive mode")
	flag.StringVar(&localitiesPathToSave, "lo", "", "[Optional] Path for result about localities. If empty, localities will not be saved")
	flag.StringVar(&villagesPathToSave, "vo", "", "[Optional] Path for result about villages. If empty, villages will not be saved")
	flag.StringVar(&categoriesPathToSave, "co", "", "[Optional] Path for result about categories. If empty, category will not be saved")
	flag.StringVar(&format, "f", "csv", "[Optional] Format for saving file, default: csv")

	flag.Parse()

	if isInteractive {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Type file name for localities:")
		localitiesPathToSave, _ = reader.ReadString('\n')
		fmt.Println("Type file name for villages:")
		villagesPathToSave, _ = reader.ReadString('\n')
		fmt.Println("Type format for files [csv, xlsx (default)]:")
		format, _ = reader.ReadString('\n')

		localitiesPathToSave = strings.TrimSpace(localitiesPathToSave)
		villagesPathToSave = strings.TrimSpace(villagesPathToSave)
		format = strings.TrimSpace(format)

		if format == "" {
			format = "xlsx"
		}
	}

	return AppParameters{
		XmlFilePath:          xmlFilePath,
		LocalitiesPathToSave: localitiesPathToSave,
		VillagesPathToSave:   villagesPathToSave,
		CategoriesPathToSave: categoriesPathToSave,
		FormatFile:           format,
		IsInteractive:        isInteractive,
	}
}

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
