package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

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
			handleCreateEvent(formData.Name, formData.Date, formData.Description, formData.DetailsID, window, s.eventServ)
		}
	}, window)
}

func handleCreateEvent(eventName string, eventDate string, eventDescription string, detailsID int, window fyne.Window, eventServ service.EventService) {
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
