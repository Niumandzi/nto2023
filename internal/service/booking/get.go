package booking

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s BookingService) GetBookings(startDate string, endDate string, eventID int, categoryName string) ([]model.BookingWithFacility, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	bookings, err := s.bookingRepo.Get(ctx, startDate, endDate, eventID, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return bookings, nil
}
