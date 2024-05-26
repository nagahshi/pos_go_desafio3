package usecase

import (
	"github.com/nagahshi/pos_go_desafio3/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute() (orders []OrderOutputDTO, err error) {
	ordersData, err := c.OrderRepository.List()
	if err != nil {
		println(err)
		println()
		return nil, err
	}

	var DTO OrderOutputDTO
	for _, order := range ordersData {
		DTO.ID = order.ID
		DTO.Price = order.Price
		DTO.Tax = order.Tax
		DTO.FinalPrice = order.FinalPrice

		orders = append(orders, DTO)
	}

	return orders, nil
}
