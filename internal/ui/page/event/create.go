package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s EventPage) CreateEvent(categoryName string, window fyne.Window, onUpdate func()) {
	formData := model.Event{}

	nameEntry := component.EntryWidget("Название")
	dateEntry := component.EntryWidget("дд.мм.гггг")
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
		if confirm {
			formData.Name = nameEntry.Text
			formData.Date = dateEntry.Text
			formData.Description = descriptionEntry.Text

			handleCreateEvent(formData, window, s.eventServ, onUpdate)
		}
	}, window)
}

func handleCreateEvent(eventData model.Event, window fyne.Window, eventServ service.EventService, onUpdate func()) {
	_, err := eventServ.CreateEvent(eventData)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие создано", "Событие успешно создано!", window)
		onUpdate()
	}
}
