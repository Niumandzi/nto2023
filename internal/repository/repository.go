package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.Event) (int, error)
	Get(ctx context.Context, categoryName string, typeName string) ([]model.EventWithDetails, error)
	Update(ctx context.Context, eventUpd model.Event) error
	Delete(ctx context.Context, eventId int) error
}

type DetailsRepository interface {
	Create(ctx context.Context, categoryName string, typeName string) (int, error)
	Get(ctx context.Context, categoryName string) ([]model.Details, error)
	GetId(ctx context.Context, categoryName string, typeName string) (int, error)
	UpdateTypeName(ctx context.Context, detailsId int, typeName string) error
	DeleteType(ctx context.Context, categoryName string, typeName string) error
}
