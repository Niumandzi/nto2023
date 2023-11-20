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

func (s DetailsPage) ShowDetails(categoryName string, window fyne.Window, eventContainer *fyne.Container) {
	details, err := s.detailsServ.GetDetails(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, detail := range details {
		card := s.createDetailCard(detail, window, func() {
			s.ShowDetails(categoryName, window, eventContainer)
		})
		grid.Add(card)
	}

	eventContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	eventContainer.Refresh()
}

func (s DetailsPage) createDetailCard(detail model.Details, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(detail)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateDetail(detail.ID, detail.Category, detail.TypeName, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.detailsServ.DeleteDetail(detail.ID)
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

func card(detail model.Details) string {
	var category string
	switch detail.Category {
	case "entertainment":
		category = "Развлечения"
	case "enlightenment":
		category = "Просвещение"
	case "education":
		category = "Образование"
	}

	return fmt.Sprintf("Категория: %s\nТип: %s",
		category, detail.TypeName)
}
