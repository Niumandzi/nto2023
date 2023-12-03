package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s EventPage) UpdateEvent(categoryName string, typeName string, event model.Event, window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWithDataWidget("Название", event.Name)
	dateEntry := component.EntryWithDataWidget("гггг-мм-дд", event.Date)
	descriptionEntry := component.MultiLineEntryWidgetWithData("Описание", event.Description)

	details, err := s.detailsServ.GetActiveDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	typeNames := make(map[string]int)
	for _, detail := range details {
		typeNames[detail.TypeName] = detail.ID
	}

	detailsSelect := component.SelectorWidget(typeName, typeNames, func(id int) {
		event.DetailsID = id
	},
		nil,
	)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", detailsSelect),
		widget.NewFormItem("", nameEntry),
		widget.NewFormItem("", dateEntry),
		widget.NewFormItem("", descriptionEntry),
	}

	dialog.ShowForm("                                Обновить событие                     ", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			event.Name = nameEntry.Text
			event.Date = dateEntry.Text
			event.Description = descriptionEntry.Text
			handleUpdateEvent(event, window, s.eventServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(event model.Event, window fyne.Window, eventServ service.EventService, onUpdate func()) {
	err := eventServ.UpdateEvent(event)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие обновлено", "Событие успешно обновлено!", window)
		onUpdate()
	}
}
