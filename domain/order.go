package domain

import (
	"github.com/joesantosio/example-order-book/entity"
)

// TODO: need to setup tests on the domain

func CancelOrder(userOrderID int, userID int, orderRepo entity.RepositoryOrder) (entity.Order, error) {
	order, err := orderRepo.GetOrderByID(userOrderID, userID)
	if err != nil {
		return entity.Order{}, err
	}

	topOrder, err := orderRepo.GetTopOrder(order.Symbol, order.Side)
	if err != nil {
		return entity.Order{}, err
	}

	_, err = orderRepo.Cancel(userOrderID, userID)
	if err != nil {
		return entity.Order{}, err
	}

	// maybe the top order was updated
	if topOrder.ID == userOrderID {
		topOrder, err = orderRepo.GetTopOrder(order.Symbol, order.Side)
		if err != nil {
			return entity.Order{}, err
		}

		if topOrder.Symbol != "" {
			return topOrder, nil
		}

		// isnt there a topOrder? we need to inform
		return entity.NewOrder(-1, -1, order.Symbol, order.Side, -1, -1, true, false, -1, -1), nil
	}

	return entity.Order{}, nil
}

func FlushAllOrders(orderRepo entity.RepositoryOrder) (bool, error) {
	return orderRepo.Flush()
}

func RequestBuy(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
) (topOrder entity.Order, isCrossBook bool, err error) {
	topSellingOrder, err := orderRepo.GetTopOrder(symbol, "ask")
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topSellingOrder.Price != 0 && limitPrice >= topSellingOrder.Price
	if isCrossBook {
		return topOrder, isCrossBook, nil
	}

	// TODO: match and trade

	_, err = orderRepo.Create(userOrderID, userID, symbol, "bid", limitPrice, qty)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	topOrder, err = orderRepo.GetTopOrder(symbol, "bid")
	if err != nil {
		return entity.Order{}, isCrossBook, err
	}

	// DEV: the count for the TopOrder is done per user
	if topOrder.UserID != userID {
		topOrder = entity.Order{}
	}

	return topOrder, isCrossBook, nil
}

func RequestSell(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
) (topOrder entity.Order, isCrossBook bool, err error) {
	topBuyingOrder, err := orderRepo.GetTopOrder(symbol, "bid")
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topBuyingOrder.Price != 0 && limitPrice <= topBuyingOrder.Price
	if isCrossBook {
		return topOrder, isCrossBook, nil
	}

	// TODO: match and trade

	_, err = orderRepo.Create(userOrderID, userID, symbol, "ask", limitPrice, qty)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	topOrder, err = orderRepo.GetTopOrder(symbol, "ask")
	if err != nil {
		return entity.Order{}, isCrossBook, err
	}

	// DEV: the count for the TopOrder is done per user
	if topOrder.UserID != userID {
		topOrder = entity.Order{}
	}

	return topOrder, isCrossBook, nil
}
