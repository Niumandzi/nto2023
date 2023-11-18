package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.EventWithDetails) (int, error)
	CreateDetails(categoryName string, typeName string) (int, error)
	GetEventsByCategory(categoryName string) ([]model.EventWithDetails, error)
	GetEventsByCategoryAndType(categoryName string, typeName string) ([]model.EventWithDetails, error)
	GetDetailsByCategory(categoryName string) ([]model.Details, error)
	UpdateEvent(eventUpd model.Event) error
	UpdateTypeName(detailsId int, typeName string) error
	DeleteEvent(eventId int) error
	DeleteType(detailsId int) error
}
