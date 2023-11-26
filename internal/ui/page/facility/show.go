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

func (s FacilityPage) createFacilityCard(facility model.Facility, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(facility)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		s.UpdateFacility(facility.ID, facility.Name, window, onUpdate)
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

func card(facility model.Facility) string {

	return fmt.Sprintf("Помещение: %s", facility.Name)
}
