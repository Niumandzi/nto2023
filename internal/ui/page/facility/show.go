package facility

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
	"strings"
)

func (s FacilityPage) ShowFacility(window fyne.Window, eventContainer *fyne.Container) {
	facility, err := s.facilityServ.GetFacilities("", 0, "")
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, facility := range facility {
		card := s.createFacilityCard(facility, window, func() {
			s.ShowFacility(window, eventContainer)
		})
		grid.Add(card)
	}

	eventContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	eventContainer.Refresh()
}

func (s FacilityPage) createFacilityCard(facility model.FacilityWithParts, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(facility)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateFacility(facility.ID, facility.Name, facility.Parts, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.facilityServ.DeleteFacility(facility.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Помещение удалено", "Помещение успешно удалено!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)
	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(facility model.FacilityWithParts) string {
	result := fmt.Sprintf("Помещение: %s", facility.Name)

	if len(facility.Parts) > 0 {
		partsInfo := []string{}
		for _, part := range facility.Parts {
			partsInfo = append(partsInfo, fmt.Sprintf("Часть: %s", part.Name))
		}
		result += "\n" + strings.Join(partsInfo, "\n")
	} else {
		result += "\nЧасти: нет"
	}

	return result
}
