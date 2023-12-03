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
	cardText, isActive := card(event)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		eventToUpdate := model.Event{
			ID:          event.ID,
			Name:        event.Name,
			Date:        event.Date,
			Description: event.Description,
			DetailsID:   event.Details.ID,
		}
		s.UpdateEvent(event.Details.Category, event.Details.TypeName, eventToUpdate, window, onUpdate)
	})

	var icon fyne.Resource
	var dialogTitle, dialogMessage string

	if isActive {
		icon = theme.DeleteIcon()
		dialogTitle = "Событие удалено"
		dialogMessage = "Событие успешно удалено!"
	} else {
		icon = theme.ViewRefreshIcon()
		dialogTitle = "Событие восстановлено"
		dialogMessage = "Событие успешно восстановлено!"
	}

	deleteButton := widget.NewButtonWithIcon("", icon, func() {
		err := s.eventServ.DeleteRestoreEvent(event.ID, !isActive)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation(dialogTitle, dialogMessage, window)
			onUpdate()
		}
	})

	//deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
	//	err := s.eventServ.DeleteRestoreEvent(event.ID, false)
	//	if err != nil {
	//		dialog.ShowError(err, window)
	//	} else {
	//		dialog.ShowInformation("Событие удалено", "Событие успешно удалено!", window)
	//		onUpdate()
	//	}
	//})

	bookingButton := widget.NewButtonWithIcon("", theme.FileIcon(), func() {

	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance
	bookingButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), bookingButton, updateButton, deleteButton)

	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(event model.EventWithDetails) (string, bool) {
	return fmt.Sprintf("Тип: %s\nНазвание: %s\nДата: %s\nОписание: %s",
		event.Details.TypeName, event.Name, event.Date, event.Description), event.IsActive
}
