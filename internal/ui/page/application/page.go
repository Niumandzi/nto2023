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
	facilityServ    service.FacilityService
	workTypeServ    service.WorkTypeService
	logger          logging.Logger
}

func NewApplicationPage(appl service.ApplicationService, fac service.FacilityService, work service.WorkTypeService, logger logging.Logger) ApplicationPage {
	return ApplicationPage{
		applicationServ: appl,
		facilityServ:    fac,
		workTypeServ:    work,
		logger:          logger,
	}
}

func (s ApplicationPage) IndexApplication(categoryName string, status string, window fyne.Window) fyne.CanvasObject {
	applicationContainer := container.NewStack()

	var selectedFacilityId int
	var selectedWorkTypeId int

	updateApplicationList := func() {
		s.ShowApplication(categoryName, selectedFacilityId, selectedWorkTypeId, status, window, applicationContainer)
	}

	facilities, err := s.facilityServ.GetFacilities()
	if err != nil {
		dialog.ShowError(err, window)
		return nil
	}

	facilityNames := make(map[string]int)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityId = id
		updateApplicationList()
	})

	workTypes, err := s.workTypeServ.GetWorkTypes()
	if err != nil {
		dialog.ShowError(err, window)
		return nil
	}

	workNames := make(map[string]int)
	for _, work := range workTypes {
		workNames[work.Name] = work.ID
	}

	workSelect := component.SelectorWidget("Тип работ", workNames, func(id int) {
		selectedWorkTypeId = id
		updateApplicationList()
	})

	sortingButtons := container.NewHBox(facilitySelect, workSelect)

	var toolbar, content *fyne.Container
	if categoryName != "" {
		createApplicationButton := widget.NewButton("Создать заявку на выполнение работ", func() {
			//s.CreateApplication(categoryName, window, func() {
			//	updateApplicationList(-1)
			//})
		})
		createButtons := container.NewHBox(createApplicationButton)
		toolbar = container.NewBorder(nil, nil, sortingButtons, createButtons)
	} else {
		toolbar = container.NewBorder(nil, nil, sortingButtons, nil)
	}
	content = container.NewBorder(toolbar, nil, nil, nil, applicationContainer)

	updateApplicationList()

	return content
}
