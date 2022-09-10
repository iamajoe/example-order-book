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

func printHelp() {
	fmt.Println("\nUsage:\n    book <input.csv> <output.csv>")
}

func flushAllOrders(responseCh chan []string, errCh chan error, record []string, repos entity.Repositories) {
	// F
	_, err := domain.FlushAllOrders(repos.GetOrder())
	if err != nil {
		errCh <- err
		return
	}

	// TODO: what to send in a success??
	responseCh <- []string{}
}

func cancelOrder(responseCh chan []string, errCh chan error, record []string, repos entity.Repositories) {
	// C, user(int),userOrderId(int)
	if len(record) < 2 {
		responseCh <- []string{}
		return
	}

	userIDRaw := record[1]
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		errCh <- err
		return
	}

	userOrderIDRaw := record[2]
	userOrderID, err := strconv.Atoi(userOrderIDRaw)
	if err != nil {
		errCh <- err
		return
	}

	err = domain.CancelOrder(userOrderID, userID, repos.GetOrder())
	if err != nil {
		errCh <- err
		return
	}

	responseCh <- []string{"A", userIDRaw, userOrderIDRaw}
}

func newOrder(responseCh chan []string, errCh chan error, record []string, repos entity.Repositories) {
	// N, user(int),symbol(string),price(int),qty(int),side(char B or S),userOrderId(int)
	if len(record) < 7 {
		responseCh <- []string{}
		return
	}

	symbol := record[2]
	side := record[5]

	userIDRaw := record[1]
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		errCh <- err
		return
	}

	price, err := strconv.Atoi(record[3])
	if err != nil {
		errCh <- err
		return
	}

	qty, err := strconv.Atoi(record[4])
	if err != nil {
		errCh <- err
		return
	}

	userOrderIDRaw := record[6]
	userOrderID, err := strconv.Atoi(userOrderIDRaw)
	if err != nil {
		errCh <- err
		return
	}

	switch side {
	case "B":
		_, err = domain.RequestBuy(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
			repos.GetStock(),
		)
		if err != nil {
			errCh <- err
			return
		}

		responseCh <- []string{"A", userIDRaw, userOrderIDRaw}
	case "S":
		_, err = domain.RequestSell(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
			repos.GetStock(),
		)
		if err != nil {
			errCh <- err
			return
		}

		responseCh <- []string{"A", userIDRaw, userOrderIDRaw}
	default:
		responseCh <- []string{}
	}
}

func Init(args []string, repos entity.Repositories) error {
	if len(args) < 3 {
		printHelp()
		return nil
	}

	records, err := readCSVFile(args[1])
	if err != nil {
		return err
	}

	responseCh := make(chan []string, len(records))
	errCh := make(chan error, 1)

	for _, record := range records {
		// ignore empty lines
		if len(record) < 1 {
			responseCh <- []string{}
			continue
		}

		switch record[0] {
		case "N":
			go newOrder(responseCh, errCh, record, repos)
		case "C":
			go cancelOrder(responseCh, errCh, record, repos)
		case "F":
			go flushAllOrders(responseCh, errCh, record, repos)
		}
	}

	// wait for all to go through...
	var resArr [][]string
	for {
		shouldBreak := false

		select {
		case res := <-responseCh:
			resArr = append(resArr, res)
			if len(resArr) == len(records) {
				shouldBreak = true
				close(responseCh)
			}
		case err := <-errCh:
			close(errCh)
			return err
		}

		if shouldBreak {
			break
		}
	}

	// ... and save
	err = writeCSVFile(args[2], resArr)

	return err
}
