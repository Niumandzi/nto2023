package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type FacilityPage struct {
	facilityServ service.FacilityService
	logger       logging.Logger
}

func NewFacilityPage(det service.FacilityService, logger logging.Logger) FacilityPage {
	return FacilityPage{
		facilityServ: det,
		logger:       logger,
	}
}

func (s FacilityPage) IndexFacility(window fyne.Window) fyne.CanvasObject {
	facilityContainer := container.NewStack()
	facilityList := func(eventType string) {
		s.ShowFacility(window, facilityContainer)
	}

	createFacilityButton := widget.NewButton("Создать новое помещение", func() {
		s.CreateFacility(window, func() {
			facilityList("")
		})
	})

	createButtons := container.NewHBox(createFacilityButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, facilityContainer)
	facilityList("")

	return content
}
