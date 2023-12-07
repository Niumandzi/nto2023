package mug

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type MugTypePage struct {
	mugTypeServ service.MugTypeService
	logger      logging.Logger
}

func NewMugTypePage(mug service.MugTypeService, logger logging.Logger) MugTypePage {
	return MugTypePage{
		mugTypeServ: mug,
		logger:      logger,
	}
}

func (m MugTypePage) IndexMugType(window fyne.Window) fyne.CanvasObject {
	teacherContainer := container.NewStack()
	teacherList := func(eventType string) {
		m.ShowMugType(window, teacherContainer)
	}

	createMugTypeButton := widget.NewButton("Создать новый тип кружка", func() {
		m.CreateMugType(window, func() {
			teacherList("")
		})
	})

	createButtons := container.NewHBox(createMugTypeButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, teacherContainer)
	teacherList("")

	return content
}
