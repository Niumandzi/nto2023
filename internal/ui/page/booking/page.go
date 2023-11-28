package booking

import (
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type BookingPage struct {
	bookingServ service.BookingService
	logger      logging.Logger
}

func NewBookingPage(book service.BookingService, logger logging.Logger) BookingPage {
	return BookingPage{
		bookingServ: book,
		logger:      logger,
	}
}
