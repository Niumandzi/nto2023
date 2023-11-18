package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s EventPage) CreateEventType(categoryName string, window fyne.Window) {
	nameEntry := component.EntryWidget("Тип события")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание нового типа события", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateDetails(nameEntry.Text, categoryName, window, s.eventServ)
		}
	}, window)
}

func handleCreateDetails(eventName string, categoryName string, window fyne.Window, eventServ service.EventService) {
	_, err := eventServ.CreateDetails(categoryName, eventName)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип создан", "Тип для события успешно создано!", window)
	}
}

func (s EventPage) CreateEvent(categoryName string, window fyne.Window) {
	formData := struct {
		Name        string
		Date        string
		Description string
		DetailsID   int
	}{}

	nameEntry := component.EntryWidget("Название")
	dateEntry := component.EntryWidget("Дата")
	descriptionEntry := component.MultiLineEntryWidget("Описание")

	details, err := s.eventServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := make(map[string]int)
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	detailsSelect := component.SelectorWidget(typeNames, func(id int) {
		formData.DetailsID = id
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
		widget.NewFormItem("", dateEntry),
		widget.NewFormItem("", descriptionEntry),
		widget.NewFormItem("", detailsSelect),
	}

	dialog.ShowForm("Создать событие", "Создать", "Отмена", formItems, func(confirm bool) {

		formData.Name = nameEntry.Text
		formData.Date = dateEntry.Text
		formData.Description = descriptionEntry.Text

		if confirm {
			handleCreateEvent(formData.Name, formData.Date, formData.Description, formData.DetailsID, categoryName, window, s.eventServ)
		}
	}, window)
}

func handleCreateEvent(eventName string, eventDate string, eventDescription string, detailsID int, categoryName string, window fyne.Window, eventServ service.EventService) {
	newEvent := model.Event{
		Name:        eventName,
		Date:        eventDate,
		Description: eventDescription,
		DetailsID:   detailsID,
	}

	_, err := eventServ.CreateEvent(newEvent)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие создано", "Событие успешно создано!", window)
	}
}
