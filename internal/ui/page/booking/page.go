package booking

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type BookingPage struct {
	bookingServ  service.BookingService
	eventServ    service.EventService
	facilityServ service.FacilityService
	logger       logging.Logger
}

func NewBookingPage(book service.BookingService, event service.EventService, fac service.FacilityService, logger logging.Logger) BookingPage {
	return BookingPage{
		bookingServ:  book,
		eventServ:    event,
		facilityServ: fac,
		logger:       logger,
	}
}

func (s BookingPage) IndexBooking(categoryName string, window fyne.Window) fyne.CanvasObject {
	bookingContainer := container.NewStack()
	bookingList := func(startDate string, endDate string, eventID int, categoryName string) {
		s.ShowBooking(startDate, endDate, eventID, categoryName, window, bookingContainer)
	}

	//details, err := s.detailsServ.GetDetails(categoryName, true)
	//if err != nil {
	//	dialog.ShowError(err, window)
	//}
	//
	typeNames := map[string]int{"Все": 0}
	//for _, detail := range details {
	//	typeNames[detail.TypeName] = detail.ID
	//}
	//
	typeSelect := component.SelectorWidget("Тип мероприятия", typeNames, func(id int) {
		bookingList("", "", 0, categoryName)
	},
		nil,
	)

	createBookingButton := widget.NewButton("Создать бронь", func() {
		s.CreateBooking(categoryName, window, func() {
			bookingList("", "", 0, categoryName)
		})
	})
	createButtons := container.NewHBox(createBookingButton)

	toolbar := container.NewBorder(nil, nil, typeSelect, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, bookingContainer)
	bookingList("", "", 0, categoryName)

	return content
}
