package details

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type DetailsPage struct {
	detailsServ service.DetailsService
	logger      logging.Logger
}

func NewDetailsPage(det service.DetailsService, logger logging.Logger) DetailsPage {
	return DetailsPage{
		detailsServ: det,
		logger:      logger,
	}
}

func (s DetailsPage) IndexDetails(categoryName string, window fyne.Window) fyne.CanvasObject {
	detailsContainer := container.NewStack()
	detailsList := func(eventType string) {
		s.ShowDetails(categoryName, window, detailsContainer)
	}

	createDetailsButton := widget.NewButton("Создать новый тип события", func() {
		s.CreateDetails(categoryName, window, func() {
			detailsList("")
		})
	})

	createButtons := container.NewHBox(createDetailsButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, detailsContainer)
	detailsList("")

	return content
}
