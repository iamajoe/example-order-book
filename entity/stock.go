package entity

type Stock struct {
	ID     int
	UserID int
	Symbol string
	Qty    int
}

func NewStock(
	id int,
	userID int,
	symbol string,
	qty int,
) Stock {
	return Stock{id, userID, symbol, qty}
}

type RepositoryStock interface {
	GetByUserAndSymbol(userID int, symbol string) (Stock, error)
	Update(
		userID int,
		symbol string,
		qty int,
	) (bool, error)
}
