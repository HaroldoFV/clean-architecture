package entity

type OrderRepositoryInterface interface {
	Save(Order *Order) error
	List() ([]*Order, error)
}
