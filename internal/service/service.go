package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.Event) (int, error)
	GetEvents(categoryName string, detailsID int) ([]model.EventWithDetails, error)
	UpdateEvent(eventUpd model.Event) error
	DeleteRestoreEvent(eventId int, isActive bool) error
}

type DetailsService interface {
	CreateDetail(categoryName string, typeName string) (int, error)
	GetDetails(categoryName string, isActive bool) ([]model.Details, error)
	UpdateDetail(detailsId int, typeName string) error
	DeleteRestoreType(detailsId int, isActive bool) error
}

type WorkTypeService interface {
	CreateWorkType(name string) (int, error)
	GetWorkTypes() ([]model.WorkType, error)
	GetActiveWorkTypes(categoryName string, facilityID int, status string) ([]model.WorkType, error)
	UpdateWorkType(workTypeId int, name string) error
	DeleteRestoreWorkType(id int, isActive bool) error
}

type FacilityService interface {
	CreateFacility(name string, parts []string) (int, error)
	GetFacilities(categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error)
	GetFacilitiesByDate(startDate string, endDate string) ([]model.FacilityWithParts, error)
	UpdateFacility(facilityId int, name string) error
	DeleteRestoreFacility(id int, isActive bool) error
}

type ApplicationService interface {
	CreateApplication(application model.Application) (int, error)
	GetApplications(categoryName string, facilityId int, workTypeId int, status string) ([]model.ApplicationWithDetails, error)
	UpdateApplication(applicationUpd model.Application) error
	DeleteApplication(applicationId int) error
}

type BookingService interface {
	CreateBooking(booking model.Booking) (int, error)
	GetBookings(startDate string, endDate string, eventID int, categoryName string) ([]model.BookingWithFacility, error)
	UpdateBooking(bookingUpd model.Booking) error
	DeleteBooking(bookingId int) error
}

type PartService interface {
	CreatePart(facilityID int, partNames []string) (int, error)
	UpdatePart(updates map[int]string) error
	DeletePart(partIds []int, isActive bool) error
}
