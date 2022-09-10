package domain

import (
	"github.com/joesantosio/example-order-book/entity"
)

// TODO: need to setup tests on the domain

func GetOrderByID(userOrderID int, userID int, orderRepo entity.RepositoryOrder) (entity.Order, error) {
	return orderRepo.GetOrderByID(userOrderID, userID)
}

func GetTopOrder(symbol string, side string, orderRepo entity.RepositoryOrder) (entity.Order, error) {
	order, err := orderRepo.GetTopOrder(symbol, side)
	if err != nil {
		return entity.Order{}, err
	}

	if order.Symbol == "" {
		return entity.NewOrder(-1, -1, symbol, side, -1, -1, true, false, -1, -1), nil
	}

	return order, nil
}

func CancelOrder(userOrderID int, userID int, orderRepo entity.RepositoryOrder) error {
	_, err := orderRepo.Cancel(userOrderID, userID)
	return err
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
) (isCrossBook bool, err error) {
	topSellingOrder, err := orderRepo.GetTopOrder(symbol, "ask")
	if err != nil {
		return isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topSellingOrder.Price != 0 && limitPrice >= topSellingOrder.Price
	if isCrossBook {
		return isCrossBook, nil
	}

	// TODO: match and trade

	_, err = orderRepo.Create(userOrderID, userID, symbol, "bid", limitPrice, qty)
	return isCrossBook, err
}

func RequestSell(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
) (isCrossBook bool, err error) {
	topBuyingOrder, err := orderRepo.GetTopOrder(symbol, "bid")
	if err != nil {
		return isCrossBook, err
	}

	// check if the request will cross book
	isCrossBook = topBuyingOrder.Price != 0 && limitPrice <= topBuyingOrder.Price
	if isCrossBook {
		return isCrossBook, nil
	}

	// TODO: match and trade

	_, err = orderRepo.Create(userOrderID, userID, symbol, "ask", limitPrice, qty)
	return isCrossBook, err
}
