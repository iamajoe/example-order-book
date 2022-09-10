package entity

type Repositories interface {
	Close() error
	GetOrder() RepositoryOrder
}
