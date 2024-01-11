package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s ApplicationPage) CreateApplication(categoryName string, window fyne.Window, onUpdate func()) {
	formData := model.Application{
		Status:     "created",
		CreateDate: time.Now().Format("2006-01-02"),
	}

	statusLabel := widget.NewLabel("Черновик")
	createDateLabel := widget.NewLabel(formData.CreateDate)
	createDateLabel.Wrapping = fyne.TextWrapWord

	descriptionEntry := component.MultiLineEntryWidget("Описание")
	dueDateEntry := component.EntryWidget("Дата выполнения (гггг-мм-дд)")

	workTypes, err := s.workTypeServ.GetActiveWorkTypes("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
	}

	workNames := make(map[string]int)
	for _, work := range workTypes {
		workNames[work.Name] = work.ID
	}

	workSelect := component.SelectorWidget("Тип работ", workNames, func(id int) {
		formData.WorkTypeId = id
	},
		nil,
	)

	facilities, err := s.facilityServ.GetActiveFacilities("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
	}

	facilityNames := make(map[string]int)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget("Помещение", facilityNames, func(id int) {
		formData.FacilityId = id
	},
		nil,
	)

	events, err := s.eventServ.GetActiveEvents(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	eventSelect := component.SelectorWidget("Мероприятие", eventNames, func(id int) {
		formData.EventId = id
	},
		nil,
	)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", statusLabel),
		widget.NewFormItem("", facilitySelect),
		widget.NewFormItem("", workSelect),
		widget.NewFormItem("", descriptionEntry),
		widget.NewFormItem("", createDateLabel),
		widget.NewFormItem("", dueDateEntry),
		widget.NewFormItem("", eventSelect),
	}

	dialog.ShowForm("                            Создать заявку на работы                           ", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			formData.Description = descriptionEntry.Text
			formData.DueDate = dueDateEntry.Text
			handleCreateApplication(formData, window, s.applicationServ, onUpdate)
		}
	}, window)
}

func handleCreateApplication(appData model.Application, window fyne.Window, applicationServ service.ApplicationService, onUpdate func()) {
	_, err := applicationServ.CreateApplication(appData)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Заявка создана", "Заявка успешно создана!", window)
		onUpdate()
	}
}
