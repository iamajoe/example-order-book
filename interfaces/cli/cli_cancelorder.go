package cli

import (
	"strconv"

	"github.com/joesantosio/example-order-book/domain"
	"github.com/joesantosio/example-order-book/entity"
)

func cancelOrder(record []string, repos entity.Repositories) ([][]string, error) {
	// C, user(int),userOrderId(int)
	response := [][]string{}

	if len(record) < 2 {
		return response, nil
	}

	userID, err := strconv.Atoi(record[1])
	if err != nil {
		return response, err
	}

	userOrderID, err := strconv.Atoi(record[2])
	if err != nil {
		return response, err
	}

	topOrder, err := domain.CancelOrder(userOrderID, userID, repos.GetOrder())
	if err != nil {
		return response, err
	}

	response = append(response, getAcknowledgeResponse(userOrderID, userID))
	if topOrder.Symbol != "" {
		response = append(response, getTopOrderResponse(topOrder))
	}

	return response, nil
}
