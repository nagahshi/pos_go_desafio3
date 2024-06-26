package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"

	"github.com/nagahshi/pos_go_desafio3/internal/infra/graph/model"
	"github.com/nagahshi/pos_go_desafio3/internal/usecase"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
	dto := usecase.OrderInputDTO{
		Price: input.Price,
		Tax:   input.Tax,
	}
	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &model.Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

// ListOrder is the resolver for the listOrder field.
func (r *queryResolver) ListOrder(ctx context.Context) ([]*model.Order, error) {
	outputList, err := r.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	var order model.Order
	for _, output := range outputList {
		order.ID = output.ID
		order.Price = output.Price
		order.Tax = output.Tax
		order.FinalPrice = output.FinalPrice

		orders = append(orders, &order)
	}

	return orders, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
