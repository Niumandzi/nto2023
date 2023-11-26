package work

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type WorkTypePage struct {
	workTypeServ service.WorkTypeService
	logger       logging.Logger
}

func NewWorkTypePage(det service.WorkTypeService, logger logging.Logger) WorkTypePage {
	return WorkTypePage{
		workTypeServ: det,
		logger:       logger,
	}
}

func (s WorkTypePage) IndexWorkType(window fyne.Window) fyne.CanvasObject {
	workTypeContainer := container.NewStack()
	workTypeList := func(eventType string) {
		s.ShowWorkType(window, workTypeContainer)
	}

	createWorkTypeButton := widget.NewButton("Создать новый тип работ", func() {
		s.CreateWorkType(window, func() {
			workTypeList("")
		})
	})

	createButtons := container.NewHBox(createWorkTypeButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, workTypeContainer)
	workTypeList("")

	return content
}
