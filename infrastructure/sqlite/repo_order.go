package sqlite

import (
	"errors"
	"fmt"
	"time"

	"github.com/joesantosio/example-order-book/entity"
)

// ------------------------------
//

func getSideValidity(side string) bool {
	return !(side != "ask" && side != "bid" && side != "buy" && side != "sell")
}

func getOpenStateBySide(side string) int {
	if side != "buy" && side != "sell" {
		return 0
	}

	return 1
}

// ------------------------------
// repository

type repositoryOrder struct {
	tableName string
	db        *DB
}

func (repo *repositoryOrder) Create(
	userOrderID int,
	userID int,
	symbol string,
	side string,
	price int,
	size int,
) (int, error) {
	timeNow := time.Now().Unix()

	// make sure the side is valid with what we expect
	if !getSideValidity(side) {
		return -1, errors.New("side is invalid")
	}

	sts, err := repo.db.db.Prepare(
		"INSERT INTO " + repo.tableName + "(id, userid, symbol, side, price, size, isopen, createdat, updatedat) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id",
	)
	if err != nil {
		return -1, err
	}
	defer sts.Close()

	isOpen := getOpenStateBySide(side)

	var id int
	err = sts.QueryRow(userOrderID, userID, symbol, side, price, size, isOpen, timeNow, timeNow).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, err
}

func (repo *repositoryOrder) CreateBuy(userOrderID int, userID int, symbol string, price int, size int) (int, error) {
	return repo.Create(userOrderID, userID, symbol, "buy", price, size)
}

func (repo *repositoryOrder) CreateSell(userOrderID int, userID int, symbol string, price int, size int) (int, error) {
	return repo.Create(userOrderID, userID, symbol, "sell", price, size)
}

func (repo *repositoryOrder) CreateRequestBuy(userOrderID int, userID int, symbol string, price int, size int) (int, error) {
	return repo.Create(userOrderID, userID, symbol, "bid", price, size)
}

func (repo *repositoryOrder) CreateRequestSell(userOrderID int, userID int, symbol string, price int, size int) (int, error) {
	return repo.Create(userOrderID, userID, symbol, "ask", price, size)
}

func (repo *repositoryOrder) GetSymbolBySide(symbol string, side string) ([]entity.Order, error) {
	orders := []entity.Order{}

	// make sure the side is valid with what we expect
	if !getSideValidity(side) {
		return orders, errors.New("side is invalid")
	}

	isOpen := getOpenStateBySide(side)
	rows, err := repo.db.db.Query(
		"SELECT id, userid, symbol, side, price, size, isopen, createdat, updatedat FROM "+repo.tableName+" WHERE symbol=$1 AND side='ask' AND isopen=$2 ORDER BY createdat ASC", symbol, isOpen,
	)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var userID int
		var symbol string
		var side string
		var price int
		var size int
		var isOpen int
		var createdAt int
		var updatedAt int

		err = rows.Scan(&id, &userID, &symbol, &side, &price, &size, &isOpen, &createdAt, &updatedAt)
		if err != nil {
			return orders, err
		}

		order := entity.NewOrder(id, userID, symbol, side, price, size, isOpen == 1, createdAt, updatedAt)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *repositoryOrder) GetSelling(symbol string) ([]entity.Order, error) {
	return repo.GetSymbolBySide(symbol, "ask")
}

func (repo *repositoryOrder) GetBuying(symbol string) ([]entity.Order, error) {
	return repo.GetSymbolBySide(symbol, "bid")
}

func (repo *repositoryOrder) Cancel(userOrderID int, userID int) (bool, error) {
	sts := "DELETE FROM " + repo.tableName + " WHERE id=? AND userid=?"
	_, err := repo.db.db.Exec(sts, userOrderID, userID)
	return true, err
}

func (repo *repositoryOrder) Empty() (bool, error) {
	sts := "DELETE FROM " + repo.tableName
	_, err := repo.db.db.Exec(sts)
	return true, err
}

func (repo *repositoryOrder) removeTable() (bool, error) {
	sts := "DROP TABLE IF EXISTS " + repo.tableName
	_, err := repo.db.db.Exec(sts)
	return true, err
}

func createRepositoryOrder(db *DB) (*repositoryOrder, error) {
	repo := repositoryOrder{"orders", db}

	// DEV: order id is been provided from outside
	sts := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
			id              INTEGER    PRIMARY KEY    NOT NULL,
			userid  			  INTEGER 									NOT NULL,
			symbol  			  TEXT											NOT NULL,
			side    			  TEXT											NOT NULL,
			price    			  INTEGER										NOT NULL,
			size    			  INTEGER										NOT NULL,
			isopen    			INTEGER										NOT NULL,
			createdat				INTEGER										NOT NULL,
			updatedat				INTEGER 									NOT NULL
		);
	`, repo.tableName)
	_, err := repo.db.db.Exec(sts)

	return &repo, err
}