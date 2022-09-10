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

	// get the current order so we can check the top order will change, we need to inform
	order, err := domain.GetOrderByID(userOrderID, userID, repos.GetOrder())
	if err != nil {
		return response, err
	}
	oldTopOrder, err := domain.GetTopOrder(order.Symbol, order.Side, repos.GetOrder())
	if err != nil {
		return response, err
	}

	// perform the task
	err = domain.CancelOrder(userOrderID, userID, repos.GetOrder())
	if err != nil {
		return response, err
	}
	response = append(response, getAcknowledgeResponse(userOrderID, userID))

	// get the new top order so we can check if it makes sense to notify
	if userOrderID == oldTopOrder.ID {
		newTopOrder, err := domain.GetTopOrder(order.Symbol, order.Side, repos.GetOrder())
		if err != nil {
			return response, err
		}

		if newTopOrder.ID != oldTopOrder.ID || newTopOrder.ID == -1 {
			response = append(response, getTopOrderResponse(newTopOrder))
		}
	}

	return response, nil
}
