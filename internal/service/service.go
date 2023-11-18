package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.Event) (int, error)
	GetEvents(categoryName string, detailsID int) ([]model.EventWithDetails, error)
	UpdateEvent(eventUpd model.Event) error
	DeleteEvent(eventId int) error

	GetDetails(categoryName string) ([]model.Details, error)
	CreateDetails(categoryName string, typeName string) (int, error)
	UpdateTypeName(detailsId int, typeName string) error
	DeleteType(detailsId int) error
}

type DetailsService interface {
}
