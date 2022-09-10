package sqlite

import (
	"fmt"

	"github.com/joesantosio/example-order-book/entity"
)

// ------------------------------
// model

// ------------------------------
// repository

type repositoryStock struct {
	tableName string
	db        *DB
}

func (repo *repositoryStock) GetByUserAndSymbol(userID int, symbol string) (entity.Stock, error) {
	var id int
	var qty int

	err := repo.db.db.QueryRow(
		"SELECT id,qty FROM "+repo.tableName+" WHERE userid=$1 AND symbol=$2", userID, symbol,
	).Scan(&id, &qty)
	if err != nil {
		return entity.Stock{}, err
	}

	return entity.NewStock(id, userID, symbol, qty), nil
}

func (repo *repositoryStock) Create(
	userID int,
	symbol string,
	qty int,
) (int, error) {
	sts, err := repo.db.db.Prepare(
		"INSERT INTO " + repo.tableName + "(userid, symbol, qty) VALUES(?, ?, ?) RETURNING id",
	)
	if err != nil {
		return -1, err
	}
	defer sts.Close()

	var id int
	err = sts.QueryRow(userID, symbol, qty).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (repo *repositoryStock) Update(
	userID int,
	symbol string,
	qty int,
) (bool, error) {
	stock, err := repo.GetByUserAndSymbol(userID, symbol)
	if err != nil {
		return false, err
	}

	// we need a row ot be able to update, so, create
	if stock.Symbol == "" {
		_, err = repo.Create(userID, symbol, qty)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	sts := "UPDATE " + repo.tableName + " SET qty=? WHERE userid=? AND symbol=?"
	_, err = repo.db.db.Exec(sts, qty, userID, symbol)
	if err != nil {
		return false, err
	}

	return true, err
}

func (repo *repositoryStock) removeTable() (bool, error) {
	sts := "DROP TABLE IF EXISTS " + repo.tableName
	_, err := repo.db.db.Exec(sts)
	return true, err
}

func createRepositoryStock(db *DB) (*repositoryStock, error) {
	repo := repositoryStock{"stocks", db}

	sts := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s(
			id              INTEGER    PRIMARY KEY    AUTOINCREMENT,
			userid  			  INTEGER 									NOT NULL,
			symbol  			  TEXT											NOT NULL,
			qty      			  INTEGER										NOT NULL
		);
	`, repo.tableName)
	_, err := repo.db.db.Exec(sts)

	return &repo, err
}
