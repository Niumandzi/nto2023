package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, contact model.Event) (int, error)
	Get(ctx context.Context, eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error)
	Update(ctx context.Context, eventInput model.Event) error
	Delete(ctx context.Context, eventId int) error
}
