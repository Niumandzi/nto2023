package booking

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type BookingRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewBookingRepository(db *sql.DB, logger logging.Logger) BookingRepository {
	return BookingRepository{
		db:     db,
		logger: logger,
	}
}

func (b BookingRepository) Create(ctx context.Context, booking model.Booking) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (b BookingRepository) Get(ctx context.Context) ([]model.BookingWithFacility, error) {
	//TODO implement me
	panic("implement me")
}

func (b BookingRepository) Update(ctx context.Context, bookingUpd model.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (b BookingRepository) Delete(ctx context.Context, bookingId int) error {
	//TODO implement me
	panic("implement me")
}
