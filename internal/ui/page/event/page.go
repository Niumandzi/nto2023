package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type EventPage struct {
	eventServ service.EventService
	logger    logging.Logger
}

func NewEventPage(event service.EventService, logger logging.Logger) EventPage {
	return EventPage{
		eventServ: event,
		logger:    logger,
	}
}

func (s EventPage) EventIndex(category string) fyne.CanvasObject {
	contactContainer := container.NewStack()

	contactTypes := map[string]string{}

	eventList := func(contactType string) {
	}

	createButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	typeSelect := component.TypeSelectWidget(contactTypes, eventList)

	toolbar := container.NewBorder(nil, nil, typeSelect, createButton)

	content := container.NewBorder(toolbar, nil, nil, nil, contactContainer)

	return content
}
