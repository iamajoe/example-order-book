package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joesantosio/example-order-book/domain"
	"github.com/joesantosio/example-order-book/entity"
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

func writeCSVFile(p string, records [][]string) error {
	f, err := os.Create(p)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating file: %v", err))
	}

	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		errors.New(fmt.Sprintf("error writing file: %v", err))
	}

	return nil
}

func flushAllOrders(record []string, repos entity.Repositories) ([][]string, error) {
	// F
	response := [][]string{}

	_, err := domain.FlushAllOrders(repos.GetOrder())
	if err != nil {
		return response, err
	}

	// DEV: we don't want the flush on the final output
	// response = append(response, []string{"F"})
	return response, nil
}

func cancelOrder(record []string, repos entity.Repositories) ([][]string, error) {
	// C, user(int),userOrderId(int)
	response := [][]string{}

	if len(record) < 2 {
		return response, nil
	}

	userIDRaw := record[1]
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return response, err
	}

	userOrderIDRaw := record[2]
	userOrderID, err := strconv.Atoi(userOrderIDRaw)
	if err != nil {
		return response, err
	}

	err = domain.CancelOrder(userOrderID, userID, repos.GetOrder())
	if err != nil {
		return response, err
	}

	response = append(response, []string{"A", userIDRaw, userOrderIDRaw})
	return response, nil
}

func newOrder(record []string, repos entity.Repositories) ([][]string, error) {
	response := [][]string{}

	// N, user(int),symbol(string),price(int),qty(int),side(char B or S),userOrderId(int)
	if len(record) < 7 {
		return response, nil
	}

	symbol := record[2]
	side := record[5]

	userIDRaw := record[1]
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return response, err
	}

	priceRaw := record[3]
	price, err := strconv.Atoi(priceRaw)
	if err != nil {
		return response, err
	}

	qtyRaw := record[4]
	qty, err := strconv.Atoi(qtyRaw)
	if err != nil {
		return response, err
	}

	userOrderIDRaw := record[6]
	userOrderID, err := strconv.Atoi(userOrderIDRaw)
	if err != nil {
		return response, err
	}

	var topOrder entity.Order
	var isCrossBook bool

	switch side {
	case "B":
		topOrder, isCrossBook, err = domain.RequestBuy(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
		)
	case "S":
		topOrder, isCrossBook, err = domain.RequestSell(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
		)
	default:
	}

	// handle the response
	if isCrossBook {
		response = append(response, []string{"R", userIDRaw, userOrderIDRaw})
	} else {
		response = append(response, []string{"A", userIDRaw, userOrderIDRaw})
		if topOrder.UserID == userID {
			response = append(response, []string{"B", side, strconv.Itoa(topOrder.Price), strconv.Itoa(topOrder.Size)})
		}
	}

	return response, nil
}

func printHelp() {
	fmt.Println("\nUsage:\n    book <input.csv> <output.csv>")
}

func Init(args []string, repos entity.Repositories) error {
	if len(args) < 3 {
		printHelp()
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
