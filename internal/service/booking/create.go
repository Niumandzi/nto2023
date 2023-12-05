package booking

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s BookingService) CreateBooking(booking model.Booking) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&booking,
		validation.Field(&booking.CreateDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&booking.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&booking.StartTime, validation.Required),
		validation.Field(&booking.EndDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&booking.EndTime, validation.Required),
		validation.Field(&booking.EventID, validation.Required),
		validation.Field(&booking.FacilityID, validation.Required),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	start, _ := time.Parse("2006-01-02 15:04", booking.StartDate+" "+booking.StartTime)
	end, _ := time.Parse("2006-01-02 15:04", booking.EndDate+" "+booking.EndTime)

	if start.After(end) {
		err := errors.New("start date and time must be earlier than or equal to end date and time")
		s.logger.Error("Date and time range error: %v", err)
		return 0, err
	}

	bookingDB := model.Booking{
		Description: booking.Description,
		CreateDate:  booking.CreateDate,
		StartDate:   booking.StartDate,
		StartTime:   booking.StartTime,
		EndDate:     booking.EndDate,
		EndTime:     booking.EndTime,
		EventID:     booking.EventID,
		FacilityID:  booking.FacilityID,
		PartIDs:     booking.PartIDs,
	}

	id, err := s.bookingRepo.Create(ctx, bookingDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}
