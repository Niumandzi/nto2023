package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
	"time"
)

type EventPage struct {
	facilityServ service.FacilityService
	bookingServ  service.BookingService
	eventServ    service.EventService
	detailsServ  service.DetailsService
	logger       logging.Logger
}

func NewEventPage(fac service.FacilityService, book service.BookingService, event service.EventService, det service.DetailsService, logger logging.Logger) EventPage {
	return EventPage{
		facilityServ: fac,
		bookingServ:  book,
		eventServ:    event,
		detailsServ:  det,
		logger:       logger,
	}
}

func (s EventPage) IndexEvent(categoryName string, window fyne.Window) fyne.CanvasObject {
	eventContainer := container.NewStack()
	eventList := func(eventType string, id int) {
		s.ShowEvent(categoryName, id, window, eventContainer)
	}

	details, err := s.detailsServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := map[string]int{"Все": 0}
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	typeSelect := component.SelectorWidget("Тип мероприятия", typeNames, func(id int) {
		eventList("", id)
	},
		nil,
	)

	createEventButton := widget.NewButton("Создать событие", func() {
		s.CreateEvent(categoryName, window, func() {
			eventList("", 0)
		})
	})

	toolbar := container.NewBorder(nil, nil, typeSelect, createEventButton)
	content := container.NewBorder(toolbar, nil, nil, nil, eventContainer)
	eventList("", 0)

	return content
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func validateDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func validateTime(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}
