package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationPage) CreateApplication(categoryName string, window fyne.Window, onUpdate func()) {
	formData := struct {
		Name        string
		Date        string
		Description string
		DetailsID   int
	}{}

	nameEntry := component.EntryWidget("Название")
	dateEntry := component.EntryWidget("дд.мм.гггг")
	descriptionEntry := component.MultiLineEntryWidget("Описание")

	details, err := s.applicationServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := make(map[string]int)
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	detailsSelect := component.SelectorWidget("Тип", typeNames, func(id int) {
		formData.DetailsID = id
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("", detailsSelect),
		widget.NewFormItem("", nameEntry),
		widget.NewFormItem("", dateEntry),
		widget.NewFormItem("", descriptionEntry),
	}

	dialog.ShowForm("                            Создать событие                           ", "Создать", "Отмена", formItems, func(confirm bool) {

		formData.Name = nameEntry.Text
		formData.Date = dateEntry.Text
		formData.Description = descriptionEntry.Text

		if confirm {
			handleCreateApplication(formData.Name, formData.Date, formData.Description, formData.DetailsID, window, s.applicationServ, onUpdate)
		}
	}, window)
}

func handleCreateApplication(applicationName string, applicationDate string, applicationDescription string, detailsID int, window fyne.Window, applicationServ service.ApplicationService, onUpdate func()) {
	newApplication := model.Application{
		Name:        applicationName,
		Date:        applicationDate,
		Description: applicationDescription,
		DetailsID:   detailsID,
	}

	_, err := applicationServ.CreateApplication(newApplication)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие создано", "Событие успешно создано!", window)
		onUpdate()
	}
}
