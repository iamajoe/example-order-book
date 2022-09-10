package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

func readCSVFile(p string) ([][]string, error) {
	var records [][]string

	f, err := os.Open(p)
	if err != nil {
		return records, errors.New(fmt.Sprintf("error reading csv: %v", err))
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	// disable field check
	csvReader.FieldsPerRecord = -1

	records, err = csvReader.ReadAll()
	if err != nil {
		return records, errors.New(fmt.Sprintf("unable to parse csv: %v", err))
	}

	// make sure all data is as expected
	for _, record := range records {
		for i, column := range record {
			record[i] = strings.TrimSpace(column)
		}
	}

	return records, nil
}
