package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationPage) UpdateApplication(categoryName string, workTypeName string, facilityName string, eventName string, application model.Application, window fyne.Window, onUpdate func()) {
	var status string

	descriptionEntry := component.EntryWithDataWidget("Описание", application.Description)
	createDateLabel := widget.NewLabel(application.CreateDate)
	dueDateEntry := component.EntryWithDataWidget("Дата выполнения (дд.мм.гггг)", application.DueDate)

	switch application.Status {
	case "created":
		status = "Черновик"
	case "todo":
		status = "К выполнению"
	case "done":
		status = "Выполнено"
	}

	statusOptions := map[string]string{"Черновик": "created", "К выполнению": "todo", "Выполнено": "done"}
	statusSelect := component.SelectorWidget(status, statusOptions,
		nil,
		func(selectedStatus string) {
			application.Status = selectedStatus
		},
	)

	workTypes, err := s.workTypeServ.GetWorkTypes("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	workNames := make(map[string]int)
	for _, work := range workTypes {
		workNames[work.Name] = work.ID
	}

	workSelect := component.SelectorWidget(workTypeName, workNames,
		func(id int) {
			application.WorkTypeId = id
		},
		nil,
	)

	facilities, err := s.facilityServ.GetFacilities("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	facilityNames := make(map[string]int)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget(facilityName, facilityNames,
		func(id int) {
			application.FacilityId = id
		},
		nil,
	)

	events, err := s.eventServ.GetEvents(categoryName, 0)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	eventSelect := component.SelectorWidget(eventName, eventNames,
		func(id int) {
			application.EventId = id
		},
		nil,
	)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", workSelect),
		widget.NewFormItem("", facilitySelect),
		widget.NewFormItem("", descriptionEntry),
		widget.NewFormItem("", createDateLabel),
		widget.NewFormItem("", dueDateEntry),
		widget.NewFormItem("", statusSelect),
		widget.NewFormItem("", eventSelect),
	}

	dialog.ShowForm("                            Обновить заявку на работы                ", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			application.Description = descriptionEntry.Text
			application.DueDate = dueDateEntry.Text
			handleUpdateApplication(application, window, s.applicationServ, onUpdate)
		}
	}, window)
}

func handleUpdateApplication(appData model.Application, window fyne.Window, applicationServ service.ApplicationService, onUpdate func()) {
	err := applicationServ.UpdateApplication(appData)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Заявка обновлена", "Заявка успешно обновлена!", window)
		onUpdate()
	}
}
