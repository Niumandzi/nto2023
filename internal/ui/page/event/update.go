package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s EventPage) UpdateEvent(categoryName string, id int, name string, date string, Description string, DetailsID int, window fyne.Window) {
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
	dateEntry := component.EntryWithDataWidget("Дата", date)
	descriptionEntry := component.MultiLineEntryWidgetWithData("Описание", Description)

	details, err := s.eventServ.GetDetails(categoryName)
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
			handleUpdateEvent(id, nameEntry.Text, dateEntry.Text, descriptionEntry.Text, formData.DetailsID, window, s.eventServ)
		}
	}, window)
}

func handleUpdateEvent(eventID int, eventName string, eventDate string, eventDescription string, detailsID int, window fyne.Window, eventServ service.EventService) {
	updatedEvent := model.Event{
		ID:          eventID,
		Name:        eventName,
		Date:        eventDate,
		Description: eventDescription,
		DetailsID:   detailsID,
	}

	err := eventServ.UpdateEvent(updatedEvent)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие обновлено", "Событие успешно обновлено!", window)
	}
}
