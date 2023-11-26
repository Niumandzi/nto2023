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
	descriptionEntry := component.EntryWithDataWidget("Описание", application.Description)
	dueDateEntry := component.EntryWithDataWidget("Дата выполнения (дд.мм.гггг)", application.DueDate)

	workTypes, err := s.workTypeServ.GetWorkTypes()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	workNames := make(map[string]int)
	for _, work := range workTypes {
		workNames[work.Name] = work.ID
	}

	workSelect := component.SelectorWidget(workTypeName, workNames, func(id int) {
		application.WorkTypeId = id
	})

	facilities, err := s.facilityServ.GetFacilities()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	facilityNames := make(map[string]int)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget(facilityName, facilityNames, func(id int) {
		application.FacilityId = id
	})

	events, err := s.eventServ.GetEvents(categoryName, -1)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	eventSelect := component.SelectorWidget(eventName, eventNames, func(id int) {
		application.EventId = id
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("", facilitySelect),
		widget.NewFormItem("", workSelect),
		widget.NewFormItem("", descriptionEntry),
		widget.NewFormItem("", dueDateEntry),
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
