package cli

import (
	"fmt"

	"github.com/joesantosio/example-order-book/entity"
)

func Init(args []string, repos entity.Repositories) error {
	if len(args) < 3 {
		fmt.Println("\nUsage:\n    book <input.csv> <output.csv>")
		return nil
	}

	fmt.Printf("\nRunning with:\n    input: '%s'\n    output '%s'\n", args[1], args[2])

	records, err := readCSVFile(args[1])
	if err != nil {
		return err
	}

	response := [][]string{}

	for _, record := range records {
		// ignore empty lines
		if len(record) < 1 {
			continue
		}

		switch record[0] {
		case "N":
			res, err := newOrder(record, repos)
			if err != nil {
				return err
			}

			response = append(response, res...)
		case "C":
			res, err := cancelOrder(record, repos)
			if err != nil {
				return err
			}

			response = append(response, res...)
		case "F":
			res, err := flushAllOrders(record, repos)
			if err != nil {
				return err
			}

			response = append(response, res...)
		}
	}

	// ... and save
	err = writeCSVFile(args[2], response)

	return err
}
