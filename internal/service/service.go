package service

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	Create(ctx context.Context, contact model.Event) (int, error)
	Get(ctx context.Context, eventArgument string) ([]model.EventWithCategoryAndType, error)
	Update(ctx context.Context, eventInput model.Event) error
	Delete(ctx context.Context, eventId int) error
}

type Services struct {
	Event EventService
}

func NewRepositories(event EventService) *Services {
	return &Services{
		Event: event,
	}
}
