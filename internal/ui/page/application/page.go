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
	eventServ       service.EventService
	facilityServ    service.FacilityService
	workTypeServ    service.WorkTypeService
	logger          logging.Logger
}

func NewApplicationPage(appl service.ApplicationService, event service.EventService, fac service.FacilityService, work service.WorkTypeService, logger logging.Logger) ApplicationPage {
	return ApplicationPage{
		applicationServ: appl,
		eventServ:       event,
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

	facilities, err := s.facilityServ.GetFacilities(categoryName, selectedWorkTypeId, status, true)
	if err != nil {
		dialog.ShowError(err, window)
	}

	facilityNames := map[string]int{"Все": 0}
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityId = id
		updateApplicationList()
	},
		nil,
	)

	workTypes, err := s.workTypeServ.GetWorkTypes(categoryName, selectedFacilityId, status)
	if err != nil {
		dialog.ShowError(err, window)
	}

	workNames := map[string]int{"Все": 0}
	for _, work := range workTypes {
		workNames[work.Name] = work.ID
	}

	workSelect := component.SelectorWidget("Тип работ", workNames, func(id int) {
		selectedWorkTypeId = id
		updateApplicationList()
	},
		nil,
	)

	sortingButtons := container.NewHBox(facilitySelect, workSelect)

	var toolbar, content *fyne.Container
	if categoryName != "" {
		statusOptions := map[string]string{"Черновик": "created", "К выполнению": "todo", "Выполнено": "done", "Все": ""}
		statusSelect := component.SelectorWidget("Статус", statusOptions, nil, func(selectedStatus string) {
			status = selectedStatus
			updateApplicationList()
		})
		sortingButtons.Add(statusSelect)
		createApplicationButton := widget.NewButton("Создать заявку на выполнение работ", func() {
			s.CreateApplication(categoryName, window, func() {
				updateApplicationList()
			})
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
