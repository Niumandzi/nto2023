package booking

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s BookingService) UpdateBooking(bookingWithFacilityUpd model.Booking) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&bookingWithFacilityUpd,
		validation.Field(&bookingWithFacilityUpd.CreateDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingWithFacilityUpd.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingWithFacilityUpd.StartTime, validation.Required),
		validation.Field(&bookingWithFacilityUpd.EndDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&bookingWithFacilityUpd.EndTime, validation.Required),
		validation.Field(&bookingWithFacilityUpd.EventID, validation.Required),
		validation.Field(&bookingWithFacilityUpd.FacilityID, validation.Required),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	start, _ := time.Parse("2006-01-02 15:04", bookingWithFacilityUpd.StartDate+" "+bookingWithFacilityUpd.StartTime)
	end, _ := time.Parse("2006-01-02 15:04", bookingWithFacilityUpd.EndDate+" "+bookingWithFacilityUpd.EndTime)

	if start.After(end) {
		err := errors.New("start date and time must be earlier than or equal to end date and time")
		s.logger.Error("Date and time range error: ", err)
		return err
	}

	bookingUpd := model.Booking{
		ID:          bookingWithFacilityUpd.ID,
		Description: bookingWithFacilityUpd.Description,
		CreateDate:  bookingWithFacilityUpd.CreateDate,
		StartDate:   bookingWithFacilityUpd.StartDate,
		StartTime:   bookingWithFacilityUpd.StartTime,
		EndDate:     bookingWithFacilityUpd.EndDate,
		EndTime:     bookingWithFacilityUpd.EndTime,
		EventID:     bookingWithFacilityUpd.EventID,
		FacilityID:  bookingWithFacilityUpd.FacilityID,
		PartIDs:     bookingWithFacilityUpd.PartIDs,
	}

	err = s.bookingRepo.Update(ctx, bookingUpd)
	if err != nil {
		return err
	}

	return nil
}
