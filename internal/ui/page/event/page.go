package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type EventPage struct {
	eventServ   service.EventService
	detailsServ service.DetailsService
	logger      logging.Logger
}

func NewEventPage(event service.EventService, det service.DetailsService, logger logging.Logger) EventPage {
	return EventPage{
		eventServ:   event,
		detailsServ: det,
		logger:      logger,
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
