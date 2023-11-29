package repository

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

type EventRepository interface {
	Create(ctx context.Context, event model.Event) (int, error)
	Get(ctx context.Context, categoryName string, detailsID int) ([]model.EventWithDetails, error)
	Update(ctx context.Context, eventUpd model.Event) error
	Delete(ctx context.Context, eventId int) error
}

type DetailsRepository interface {
	Create(ctx context.Context, categoryName string, typeName string) (int, error)
	Get(ctx context.Context, categoryName string) ([]model.Details, error)
	GetId(ctx context.Context, categoryName string, typeName string) (int, error)
	UpdateTypeName(ctx context.Context, detailsId int, typeName string) error
	DeleteType(ctx context.Context, detailsId int) error
}

type WorkTypeRepository interface {
	Create(ctx context.Context, name string) (int, error)
	Get(ctx context.Context, categoryName string, facilityID int, status string) ([]model.WorkType, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, id int) error
}

type FacilityRepository interface {
	Create(ctx context.Context, name string, parts []string) (int, error)
	Get(ctx context.Context, categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error)
	Update(ctx context.Context, idOld int, nameUpd string) error
	Delete(ctx context.Context, id int) error
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
