package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	List() (orders []Order, err error)
	// GetTotal() (int, error)
}
