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
	GetWorkTypes() ([]model.WorkType, error)
	UpdateWorkType(workTypeId int, name string) error
	DeleteWorkType(id int) error
}

type FacilityService interface {
	CreateFacility(name string) (int, error)
	GetFacilities() ([]model.Facility, error)
	UpdateFacility(facilityId int, name string) error
	DeleteFacility(id int) error
}

type ApplicationService interface {
	CreateApplication(application model.Application) (int, error)
	GetApplications(categoryName string, facilityId int, workTypeId int, status string) ([]model.ApplicationWithDetails, error)
	UpdateApplication(applicationUpd model.Application) error
	DeleteApplication(applicationId int) error
}
