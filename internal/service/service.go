package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.Event) (int, error)
	GetEvents(categoryName string, detailsID int) ([]model.EventWithDetails, error)
	UpdateEvent(eventUpd model.Event) error
	DeleteEvent(eventId int) error
	GetDetails(categoryName string) ([]model.Details, error) //переписать
}

type DetailsService interface {
	CreateDetail(categoryName string, typeName string) (int, error)
	GetDetails(categoryName string) ([]model.Details, error)
	UpdateDetail(detailsId int, typeName string) error
	DeleteDetail(detailsId int) error
}

type WorkTypeService interface {
	CreateWorkType(name string) (int, error)
	GetAllWorkTypes() ([]model.WorkType, error)
	UpdateWorkType(workType model.WorkType) error
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
	GetApplications(categoryName string, workType string, status string) ([]model.ApplicationWithDetails, error)
	UpdateApplication(applicationUpd model.Application) error
	DeleteApplication(applicationId int) error
}
