package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type ApplicationPage struct {
	applicationServ service.ApplicationService
	logger          logging.Logger
}

func NewApplicationPage(application service.ApplicationService, logger logging.Logger) ApplicationPage {
	return ApplicationPage{
		applicationServ: application,
		logger:          logger,
	}
}

func (s ApplicationPage) IndexApplication(categoryName string, window fyne.Window) fyne.CanvasObject {
	applicationContainer := container.NewStack()
	applicationList := func(applicationType string, id int) {
		s.ShowApplication(categoryName, id, window, applicationContainer)
	}

	details, err := s.applicationServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := make(map[string]int)
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	typeSelect := component.SelectorWidget("Тип мероприятия", typeNames, func(id int) {
		applicationList("", id)
	})

	createApplicationButton := widget.NewButton("Создать событие", func() {
		s.CreateApplication(categoryName, window, func() {
			applicationList("", -1)
		})
	})
	createButtons := container.NewHBox(createApplicationButton)

	toolbar := container.NewBorder(nil, nil, typeSelect, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, applicationContainer)
	applicationList("", -1)

	return content
}
