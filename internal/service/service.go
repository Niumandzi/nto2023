package service

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(ctx context.Context, contact model.Event) (int, error)
	GetEvent(ctx context.Context, eventArgument string) ([]model.EventWithCategoryAndType, error)
	UpdateEvent(ctx context.Context, eventInput model.Event) error
	DeleteEvent(ctx context.Context, eventId int) error
}
