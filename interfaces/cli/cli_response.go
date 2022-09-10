package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joesantosio/example-order-book/entity"
)

func getRejectResponse(userOrderID int, userID int) []string {
	return []string{"R", strconv.Itoa(userID), strconv.Itoa(userOrderID)}
}

func getAcknowledgeResponse(userOrderID int, userID int) []string {
	return []string{"A", strconv.Itoa(userID), strconv.Itoa(userOrderID)}
}

func getTopOrderResponse(order entity.Order) []string {
	side := order.Side
	if order.Side == "ask" {
		side = "S"
	} else if order.Side == "bid" {
		side = "B"
	}

	priceRes := "-"
	sizeRes := "-"
	if order.ID > -1 {
		priceRes = strconv.Itoa(order.Price)
		sizeRes = strconv.Itoa(order.Size)
	}

	return []string{"B", side, priceRes, sizeRes}
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
