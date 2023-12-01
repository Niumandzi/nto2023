package booking

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (s BookingService) CreateBooking(booking model.Booking) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&booking)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	bookingDB := model.Booking{
		ID:          booking.ID,
		Description: booking.Description,
		CreateDate:  booking.CreateDate,
		StartDate:   booking.StartDate,
		EndDate:     booking.EndDate,
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
