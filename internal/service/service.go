package service

import (
	"github.com/niumandzi/nto2023/model"
)

type EventService interface {
	CreateEvent(event model.Event) (int, error)
	GetEvents(categoryName string, detailsID int, isActive bool) ([]model.EventWithDetails, error)
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
	GetWorkTypes(categoryName string, facilityID int, status string, isActive bool) ([]model.WorkType, error)
	UpdateWorkType(workTypeId int, name string) error
	DeleteRestoreWorkType(id int, isActive bool) error
}

type FacilityService interface {
	CreateFacility(name string, parts []string) (int, error)
	GetFacilities(categoryName string, workTypeID int, status string, isActive bool) ([]model.FacilityWithParts, error)
	GetFacilitiesByDate(startDate string, endDate string, isActive bool) ([]model.FacilityWithParts, error)
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
	GetBookings(categoryName string) ([]model.BookingWithFacility, error)
	UpdateBooking(bookingUpd model.Booking) error
	DeleteBooking(bookingId int) error
}
