package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

type EventForm struct {
	Event model.Event
	Types map[string]int
}

func (s EventPage) CreateEvent(categoryName string, window fyne.Window, onUpdate func()) {
	var form EventForm

	details, err := s.eventServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	form.Types = make(map[string]int)
	for _, detail := range details {
		form.Types[detail.TypeName] = detail.ID
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Тип", createTypeSelector(&form)),
		widget.NewFormItem("Название", component.EntryWidget("Название")),
		widget.NewFormItem("Дата", component.EntryWidget("дд.мм.гггг")),
		widget.NewFormItem("Описание", component.MultiLineEntryWidget("Описание")),
	}

	dialog.ShowForm("Создать событие", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateEvent(form.Event, window, s.eventServ, onUpdate)
		}
	}, window)
}

func createTypeSelector(form *EventForm) *widget.Select {
	typeNames := []string{}

	for name := range form.Types {
		typeNames = append(typeNames, name)
	}

	selectWidget := widget.NewSelect(typeNames, func(selected string) {
		form.Event.DetailsID = form.Types[selected]
	})

	return selectWidget
}
func handleCreateEvent(event model.Event, window fyne.Window, eventServ service.EventService, onUpdate func()) {
	_, err := eventServ.CreateEvent(event)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	dialog.ShowInformation("Событие создано", "Событие успешно создано!", window)
	onUpdate()
}
