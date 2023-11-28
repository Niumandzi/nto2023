package booking

import (
	"fyne.io/fyne/v2"
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

func (s BookingPage) IndexBooking(categoryName string, window fyne.Window) fyne.CanvasObject {
	//TODO implement me
	panic("implement me")
}
