package work

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

func (s WorkTypePage) ShowWorkType(window fyne.Window, eventContainer *fyne.Container) {
	workType, err := s.workTypeServ.GetWorkTypes("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, workType := range workType {
		card := s.createWorkTypeCard(workType, window, func() {
			s.ShowWorkType(window, eventContainer)
		})
		grid.Add(card)
	}

	eventContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	eventContainer.Refresh()
}

func (s WorkTypePage) createWorkTypeCard(workType model.WorkType, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(workType)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateWorkType(workType.ID, workType.Name, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.workTypeServ.DeleteWorkType(workType.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Тип удален", "Тип успешно удален!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)
	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(workType model.WorkType) string {
	return fmt.Sprintf("Тип: %s", workType.Name)
}
