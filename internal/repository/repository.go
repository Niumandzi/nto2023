package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.Event) (int, error)
	Get(ctx context.Context, categoryName string, detailsID int, isActive bool) ([]model.EventWithDetails, error)
	Update(ctx context.Context, eventUpd model.Event) error
	Delete(ctx context.Context, eventId int, isActive bool) error
}

type DetailsRepository interface {
	Create(ctx context.Context, categoryName string, typeName string) (int, error)
	Get(ctx context.Context, categoryName string, isActive bool) ([]model.Details, error)
	GetId(ctx context.Context, categoryName string, typeName string) (int, error)
	UpdateTypeName(ctx context.Context, detailsId int, typeName string) error
	DeleteRestoreType(ctx context.Context, detailsId int, isActive bool) error
}

type WorkTypeRepository interface {
	Create(ctx context.Context, name string) (int, error)
	Get(ctx context.Context, categoryName string, facilityID int, status string, isActive bool) ([]model.WorkType, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, workTypeId int, isActive bool) error
}

type FacilityRepository interface {
	Create(ctx context.Context, name string, parts []string) (int, error)
	Get(ctx context.Context, categoryName string, workTypeID int, status string, isActive bool) ([]model.FacilityWithParts, error)
	GetByDate(ctx context.Context, startDate string, endDate string, isActive bool) ([]model.FacilityWithParts, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, facilityId int, isActive bool) error
}

type ApplicationRepository interface {
	Create(ctx context.Context, application model.Application) (int, error)
	Get(ctx context.Context, categoryName string, facilityId int, workTypeId int, status string) ([]model.ApplicationWithDetails, error)
	Update(ctx context.Context, applicationUpd model.Application) error
	Delete(ctx context.Context, applicationId int) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking model.Booking) (int, error)
	Get(ctx context.Context, categoryName string) ([]model.BookingWithFacility, error)
	Update(ctx context.Context, bookingUpd model.Booking) error
	Delete(ctx context.Context, bookingId int) error
}
