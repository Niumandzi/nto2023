package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	GetEvents(categoryName string, detailsID int) ([]model.EventWithDetails, error)
	CreateEvent(event model.Event) (int, error)
	UpdateEvent(eventUpd model.Event) error
	DeleteEvent(eventId int) error
	GetDetails(categoryName string) ([]model.Details, error) //переписать
}

type DetailsService interface {
	GetDetails(categoryName string) ([]model.Details, error)
	CreateDetail(categoryName string, typeName string) (int, error)
	UpdateDetail(detailsId int, typeName string) error
	DeleteDetail(detailsId int) error
}
