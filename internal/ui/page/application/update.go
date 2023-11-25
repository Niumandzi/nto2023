package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationPage) UpdateApplication(categoryName string, id int, name string, date string, Description string, DetailsID int, window fyne.Window, onUpdate func()) {
	formData := struct {
		Name        string
		Date        string
		Description string
		DetailsID   int
	}{
		Name:        name,
		Date:        date,
		Description: Description,
		DetailsID:   DetailsID,
	}

	nameEntry := component.EntryWithDataWidget("Название", name)
	dateEntry := component.EntryWithDataWidget("дд.мм.гггг", date)
	descriptionEntry := component.MultiLineEntryWidgetWithData("Описание", Description)

	details, err := s.applicationServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := make(map[string]int)
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	detailsSelect := component.SelectorWidget("Тип мероприятия", typeNames, func(id int) {
		formData.DetailsID = id
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("", detailsSelect),
		widget.NewFormItem("", nameEntry),
		widget.NewFormItem("", dateEntry),
		widget.NewFormItem("", descriptionEntry),
	}

	dialog.ShowForm("                                Обновить событие                     ", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateApplication(id, nameEntry.Text, dateEntry.Text, descriptionEntry.Text, formData.DetailsID, window, s.applicationServ, onUpdate)
		}
	}, window)
}

func handleUpdateApplication(applicationID int, applicationName string, applicationDate string, applicationDescription string, detailsID int, window fyne.Window, applicationServ service.ApplicationService, onUpdate func()) {
	updatedApplication := model.Application{
		ID:          applicationID,
		Name:        applicationName,
		Date:        applicationDate,
		Description: applicationDescription,
		DetailsID:   detailsID,
	}

	err := applicationServ.UpdateApplication(updatedApplication)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие обновлено", "Событие успешно обновлено!", window)
		onUpdate()
	}
}
