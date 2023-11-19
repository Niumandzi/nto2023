package event

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
)

func (s EventPage) ShowEvent(categoryName string, detailsID int, window fyne.Window, eventContainer *fyne.Container) {
	events, err := s.eventServ.GetEvents(categoryName, detailsID)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, event := range events {
		card := s.createEventCard(event, window, func() {
			s.ShowEvent(categoryName, detailsID, window, eventContainer)
		})
		grid.Add(card)
	}

	eventContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	eventContainer.Refresh()
}

func (s EventPage) createEventCard(event model.EventWithDetails, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(event)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateEvent(event.Details.Category, event.ID, event.Name, event.Date, event.Description, event.Details.ID, window)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.eventServ.DeleteEvent(event.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Событие удалено", "Событие успешно удалено!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)
	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(event model.EventWithDetails) string {
	return fmt.Sprintf("Тип: %s\nНазвание: %s\nДата: %s\nОписание: %s",
		event.Details.TypeName, event.Name, event.Date, event.Description)
}
