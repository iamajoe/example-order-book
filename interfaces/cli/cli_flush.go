package cli

import (
	"github.com/joesantosio/example-order-book/domain"
	"github.com/joesantosio/example-order-book/entity"
)

func flushAllOrders(record []string, repos entity.Repositories) ([][]string, error) {
	// F
	response := [][]string{}

	_, err := domain.FlushAllOrders(repos.GetOrder())
	if err != nil {
		return response, err
	}

	return response, nil
}
