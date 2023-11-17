package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.EventWithCategoryAndType) (int, error)
	CreateType(typeName string, categoryName string) (int, error)
	GetEventsByCategory(eventCategory string) ([]model.EventWithCategoryAndType, error)
	GetEventsByCategoryAndType(eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error)
	GetTypesByCategory(eventCategory string) ([]model.EventType, error)
	UpdateEvent(eventUpd model.Event) error
	UpdateTypeName(eventCategory string, eventType string) error
	DeleteEvent(eventId int) error
	DeleteType(eventType string) error
}
