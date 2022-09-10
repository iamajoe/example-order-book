package domain

import (
	"github.com/joesantosio/example-order-book/entity"
)

// TODO: need to setup tests on the domain

func CancelOrder(userOrderID int, userID int, orderRepo entity.RepositoryOrder) error {
	_, err := orderRepo.Cancel(userOrderID, userID)
	return err
}

func FlushAllOrders(orderRepo entity.RepositoryOrder) (bool, error) {
	return orderRepo.Empty()
}

func RequestBuy(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
) (topOrder entity.Order, isCrossBook bool, err error) {
	topSellingOrder, err := orderRepo.GetSellingTopOrder(symbol)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topSellingOrder.Price != 0 && limitPrice >= topSellingOrder.Price
	if isCrossBook {
		return topOrder, isCrossBook, nil
	}

	_, err = orderRepo.CreateRequestBuy(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// TODO: match and trade

	topOrder, err = orderRepo.GetBuyingTopOrder(symbol)
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
	topBuyingOrder, err := orderRepo.GetBuyingTopOrder(symbol)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topBuyingOrder.Price != 0 && limitPrice <= topBuyingOrder.Price
	if isCrossBook {
		return topOrder, isCrossBook, nil
	}

	_, err = orderRepo.CreateRequestSell(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return topOrder, isCrossBook, err
	}

	// TODO: match and trade

	topOrder, err = orderRepo.GetSellingTopOrder(symbol)
	if err != nil {
		return entity.Order{}, isCrossBook, err
	}

	// DEV: the count for the TopOrder is done per user
	if topOrder.UserID != userID {
		topOrder = entity.Order{}
	}

	return topOrder, isCrossBook, nil
}
