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
	facility, err := s.facilityServ.GetFacilities()
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
	cardText, isActive := card(facility)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateFacility(facility.ID, facility.Name, facility.Parts, window, onUpdate)
	})

	var icon fyne.Resource
	var dialogTitle, dialogMessage string

	if isActive {
		icon = theme.DeleteIcon()
		dialogTitle = "Помещение удалено"
		dialogMessage = "Помещение успешно удален!"
	} else {
		icon = theme.ViewRefreshIcon()
		dialogTitle = "Помещение восстановлено"
		dialogMessage = "Помещение успешно восстановлен!"
	}

	deleteButton := widget.NewButtonWithIcon("", icon, func() {
		err := s.facilityServ.DeleteRestoreFacility(facility.ID, !isActive)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation(dialogTitle, dialogMessage, window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)
	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(facility model.FacilityWithParts) (string, bool) {
	result := fmt.Sprintf("Помещение: %s", facility.Name)

	if len(facility.Parts) > 0 {
		partsInfo := []string{}
		for _, part := range facility.Parts {
			partDisplay := part.Name
			if !part.IsActive {
				partDisplay += " (удалено)"
			}
			partsInfo = append(partsInfo, fmt.Sprintf("Часть: %s", partDisplay))
		}
		result += "\n" + strings.Join(partsInfo, "\n")
	} else {
		result += "\nЧасти: нет"
	}

	return result, facility.IsActive
}
