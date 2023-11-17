package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.EventWithCategoryAndType) (int, error)
	Get(ctx context.Context, eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error)
	Update(ctx context.Context, eventUpd model.Event) error
	Delete(ctx context.Context, eventId int) error
}

type CategoryTypeRepository interface {
	CreateCategoryWithType(ctx context.Context, eventCategory string, eventType string) (int, error)
	GetCategoryTypes(ctx context.Context, eventCategory string) ([]model.EventType, error)
	GetCategoryTypeID(ctx context.Context, eventCategory string, eventType string) (int, error)
	UpdateTypeName(ctx context.Context, eventTypeID int, eventType string) error
	DeleteType(ctx context.Context, eventType string) error
}
