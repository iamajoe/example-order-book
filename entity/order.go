package entity

type Order struct {
	ID         int
	UserID     int
	Symbol     string
	Side       string // BUY / SELL | ASK / BID
	Price      int
	Size       int
	IsOpen     bool
	IsCanceled bool
	CreatedAt  int
	UpdatedAt  int
}

func NewOrder(
	id int,
	userID int,
	symbol string,
	side string,
	price int,
	size int,
	isOpen bool,
	isCanceled bool,
	createdAt int,
	updatedAt int,
) Order {
	return Order{id, userID, symbol, side, price, size, isOpen, isCanceled, createdAt, updatedAt}
}

type RepositoryOrder interface {
	Create(userOrderID int, userID int, symbol string, side string, price int, size int) (int, error)
	GetOrderByID(userOrderID int, userID int) (Order, error)
	GetTopOrder(symbol string, side string) (Order, error)
	Cancel(userOrderID int, userID int) (bool, error)
	Flush() (bool, error)
}
