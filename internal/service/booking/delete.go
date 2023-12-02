package booking

import "context"

func (s BookingService) DeleteBooking(bookingID int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.bookingRepo.Delete(ctx, bookingID)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
