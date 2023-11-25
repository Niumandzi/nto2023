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

type WorkTypeService interface {
	CreateWorkType(name string) (int, error)
	GetAllWorkTypes() ([]model.WorkType, error)
	UpdateWorkType(workTypeId int, name string) error
	DeleteWorkType(id int) error
}

type FacilityService interface {
	CreateFacility(name string) (int, error)
	GetAllFacilities() ([]model.Facility, error)
	UpdateFacility(facility model.Facility) error
	DeleteFacility(id int) error
}

type ApplicationService interface {
	CreateApplication(application model.Application) (int, error)
	GetApplicationsBy(workType string, status string) ([]model.ApplicationWithDetails, error)
	UpdateApplication(applicationUpd model.Application) error
	DeleteApplication(applicationId int) error
}
