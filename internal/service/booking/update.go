package booking

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s BookingService) UpdateBooking(bookingUpd model.Booking) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&bookingUpd,
		validation.Field(&bookingUpd.CreateDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingUpd.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingUpd.StartTime, validation.Required),
		validation.Field(&bookingUpd.EndDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingUpd.EndTime, validation.Required),
		validation.Field(&bookingUpd.EventID, validation.Required),
		validation.Field(&bookingUpd.FacilityID, validation.Required),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	start, _ := time.Parse("2006-01-02 15:04", bookingUpd.StartDate+" "+bookingUpd.StartTime)
	end, _ := time.Parse("2006-01-02 15:04", bookingUpd.EndDate+" "+bookingUpd.EndTime)

	if start.After(end) {
		err := errors.New("start date and time must be earlier than or equal to end date and time")
		s.logger.Error("Date and time range error: ", err)
		return err
	}

	err = s.bookingRepo.Update(ctx, bookingUpd)
	if err != nil {
		return err
	}

	return nil
}
