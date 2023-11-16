package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(contact model.Event) (int, error)
	GetEvent(eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error)
	UpdateEvent(eventInput model.Event) error
	DeleteEvent(eventId int) error
}
