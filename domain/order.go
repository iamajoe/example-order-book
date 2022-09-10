package domain

import (
	"fmt"

	"github.com/joesantosio/example-order-book/entity"
)

func CancelOrder(userOrderID int, userID int, orderRepo entity.RepositoryOrder) error {
	_, err := orderRepo.Cancel(userOrderID, userID)
	return err
}

func FlushAllOrders(orderRepo entity.RepositoryOrder) (bool, error) {
	return orderRepo.Empty()
}

func setTransactionsOfSymbol(
	symbol string,
	orderRepo entity.RepositoryOrder,
	stockRepo entity.RepositoryStock,
) error {
	// figure if anyone is buying
	buying, err := orderRepo.GetBuying(symbol)
	if err != nil {
		return err
	}

	// nothing to do if no one is buying
	if len(buying) == 0 {
		return nil
	}

	// figure if anyone is selling
	selling, err := orderRepo.GetSelling(symbol)
	if err != nil {
		return err
	}

	// nothing to do if no one is selling
	if len(selling) == 0 {
		return nil
	}

	// TODO: buy at the lowest price
	// TODO: sell at the highest price

	return nil
}

func RequestBuy(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
	stockRepo entity.RepositoryStock,
) (int, error) {
	_, err := orderRepo.CreateRequestBuy(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return -1, err
	}

	go func() {
		err := setTransactionsOfSymbol(symbol, orderRepo, stockRepo)
		// TODO: what about errors? need also to handle the buying success
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}()

	return userOrderID, err
}

func RequestSell(
	userOrderID int,
	userID int,
	symbol string,
	limitPrice int,
	qty int,
	orderRepo entity.RepositoryOrder,
	stockRepo entity.RepositoryStock,
) (int, error) {
	_, err := orderRepo.CreateRequestSell(userOrderID, userID, symbol, limitPrice, qty)
	if err != nil {
		return -1, err
	}

	go func() {
		err := setTransactionsOfSymbol(symbol, orderRepo, stockRepo)
		// TODO: what about errors? need also to handle the selling success
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}()

	return userOrderID, err
}
