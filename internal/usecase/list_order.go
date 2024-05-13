package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func (c *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {
	return nil, nil
}
