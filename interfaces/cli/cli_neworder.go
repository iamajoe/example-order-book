package cli

import (
	"strconv"

	"github.com/joesantosio/example-order-book/domain"
	"github.com/joesantosio/example-order-book/entity"
)

func newOrder(record []string, repos entity.Repositories) ([][]string, error) {
	response := [][]string{}

	// N, user(int),symbol(string),price(int),qty(int),side(char B or S),userOrderId(int)
	if len(record) < 7 {
		return response, nil
	}

	symbol := record[2]
	side := record[5]

	userID, err := strconv.Atoi(record[1])
	if err != nil {
		return response, err
	}

	price, err := strconv.Atoi(record[3])
	if err != nil {
		return response, err
	}

	qty, err := strconv.Atoi(record[4])
	if err != nil {
		return response, err
	}

	userOrderID, err := strconv.Atoi(record[6])
	if err != nil {
		return response, err
	}

	var isCrossBook bool
	var oldTopOrder entity.Order
	var newTopOrder entity.Order

	switch side {
	case "B":
		oldTopOrder, err = domain.GetTopOrder(symbol, "bid", repos.GetOrder())
		if err != nil {
			return response, err
		}

		isCrossBook, err = domain.RequestBuy(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
		)

		newTopOrder, err = domain.GetTopOrder(symbol, "bid", repos.GetOrder())
		if err != nil {
			return response, err
		}
	case "S":
		oldTopOrder, err = domain.GetTopOrder(symbol, "ask", repos.GetOrder())
		if err != nil {
			return response, err
		}

		isCrossBook, err = domain.RequestSell(
			userOrderID,
			userID,
			symbol,
			price,
			qty,
			repos.GetOrder(),
		)

		newTopOrder, err = domain.GetTopOrder(symbol, "ask", repos.GetOrder())
		if err != nil {
			return response, err
		}
	default:
	}

	if err != nil {
		return response, err
	}

	// handle the response
	if isCrossBook {
		response = append(response, getRejectResponse(userOrderID, userID))
	} else {
		response = append(response, getAcknowledgeResponse(userOrderID, userID))
		if newTopOrder.ID != oldTopOrder.ID || newTopOrder.ID == -1 {
			response = append(response, getTopOrderResponse(newTopOrder))
		}
	}

	return response, nil
}
