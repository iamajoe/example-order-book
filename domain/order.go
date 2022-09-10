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
) (entity.Order, error) {
	id, err := orderRepo.CreateRequestBuy(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return entity.Order{}, err
	}

	// TODO: still missing stuff
	// TODO: what is crossed?! need that for rejections

	topOrder, err := orderRepo.GetBuyingTopOrder(symbol)
	if err != nil {
		return entity.Order{}, err
	}

	if topOrder.ID != id {
		return entity.Order{}, nil
	}

	return topOrder, nil
}

func RequestSell(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
) (entity.Order, error) {
	id, err := orderRepo.CreateRequestSell(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return entity.Order{}, err
	}

	// TODO: still missing stuff
	// TODO: what is crossed?! need that for rejections

	topOrder, err := orderRepo.GetSellingTopOrder(symbol)
	if err != nil {
		return entity.Order{}, err
	}

	if topOrder.ID != id {
		return entity.Order{}, nil
	}

	return topOrder, nil
}
