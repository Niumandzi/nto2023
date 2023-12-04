package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.Event) (int, error)
	Get(ctx context.Context, categoryName string, detailsID int) ([]model.EventWithDetails, error)
	GetActive(ctx context.Context, categoryName string) ([]model.EventWithDetails, error)
	Update(ctx context.Context, eventUpd model.Event) error
	Delete(ctx context.Context, eventID int, isActive bool) error
}

type DetailsRepository interface {
	Create(ctx context.Context, categoryName string, typeName string) (int, error)
	Get(ctx context.Context, categoryName string) ([]model.Details, error)
	//GetId(ctx context.Context, categoryName string, typeName string) (int, error)
	GetActive(ctx context.Context, categoryName string) ([]model.Details, error)
	Update(ctx context.Context, detailsId int, typeName string) error
	Delete(ctx context.Context, detailsID int, isActive bool) error
}

type WorkTypeRepository interface {
	Create(ctx context.Context, name string) (int, error)
	Get(ctx context.Context) ([]model.WorkType, error)
	GetActive(ctx context.Context, categoryName string, facilityID int, status string) ([]model.WorkType, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, workTypeID int, isActive bool) error
}

type FacilityRepository interface {
	Create(ctx context.Context, name string, parts []string) (int, error)
	Get(ctx context.Context) ([]model.FacilityWithParts, error)
	GetActive(ctx context.Context, categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error)
	GetByDate(ctx context.Context, startDate string, endDate string) ([]model.FacilityWithParts, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, facilityID int, isActive bool) error
}

type ApplicationRepository interface {
	Create(ctx context.Context, application model.Application) (int, error)
	Get(ctx context.Context, categoryName string, facilityID int, workTypeID int, status string) ([]model.ApplicationWithDetails, error)
	Update(ctx context.Context, applicationUpd model.Application) error
	Delete(ctx context.Context, applicationID int) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking model.Booking) (int, error)
	Get(ctx context.Context, startDate string, endDate string, eventID int, categoryName string) ([]model.BookingWithFacility, error)
	Update(ctx context.Context, bookingUpd model.Booking) error
	Delete(ctx context.Context, bookingID int) error
}

type PartRepository interface {
	Create(ctx context.Context, facilityID int, partNames []string) (int, error)
	Update(ctx context.Context, update map[int]string) error
	Delete(ctx context.Context, delete map[int]bool) error
}

// get parts - get partIds
type BookingPartRepository interface {
	Create(ctx context.Context, bookingId int, partId int) error
	GetParts(ctx context.Context, bookingId int) ([]int, error)
	Delete(ctx context.Context, bookingId int, partId int) error
}
