package sqlite

import (
	"errors"

	"github.com/joesantosio/example-order-book/entity"
)

type repositories struct {
	db    *DB
	order entity.RepositoryOrder
	stock entity.RepositoryStock
}

func (r *repositories) GetOrder() entity.RepositoryOrder {
	return r.order
}

func (r *repositories) GetStock() entity.RepositoryStock {
	return r.stock
}

func (r *repositories) Close() error {
	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			return err
		}

		r.db = nil
	}

	return nil
}

func InitRepos(db *DB) (repos entity.Repositories, err error) {
	if db == nil {
		return repos, errors.New("database didn't came in")
	}

	order, err := createRepositoryOrder(db)
	if err != nil {
		return repos, err
	}

	stock, err := createRepositoryStock(db)
	if err != nil {
		return repos, err
	}

	return &repositories{db, order, stock}, nil
}
