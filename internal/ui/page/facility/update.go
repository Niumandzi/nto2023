package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityPage) UpdateFacility(id int, name string, parts []model.Part, window fyne.Window, onUpdate func()) {
	partsEntries := make([]*widget.Entry, 0)
	updateData := make(map[int]string)

	vbox := container.NewVBox()
	nameEntry := component.EntryWidgetWithData("Помещение", name)
	vbox.Add(nameEntry)

	partsVBox := container.NewVBox()
	vbox.Add(partsVBox)

	partCount := 0
	for _, part := range parts {
		if partCount >= 2 {
			break
		}
		partCount++
		partID := part.ID

		partEntry := component.EntryWidgetWithData("Часть помещения", part.Name)
		partEntry.Resize(fyne.NewSize(500, 36))

		partEntry.OnChanged = func(newText string) {
			if newText != part.Name && newText != "" {
				updateData[partID] = newText
			}
		}

		partsVBox.Add(partEntry)
	}

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Обновить            ", func() {
		var createData []string
		for _, entry := range partsEntries {
			createData = append(createData, entry.Text)
		}
		handleUpdateFacility(id, nameEntry.Text, updateData, createData, window, s.facilityServ, s.partServ, onUpdate, customPopUp)
	})

	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)
	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 100))
	customPopUp.Show()
}

func handleUpdateFacility(id int, name string, update map[int]string, create []string, window fyne.Window, facilityServ service.FacilityService, partServ service.PartService, onUpdate func(), popUp *widget.PopUp) {
	err := facilityServ.UpdateFacility(id, name)

	if len(update) > 0 {
		err = partServ.UpdatePart(update)
	}

	if len(create) > 0 {
		_, err = partServ.CreatePart(id, create)
	}

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		popUp.Hide()
		dialog.ShowInformation("Помещение обновлено", "Помещение успешно обновлено!", window)
		onUpdate()
	}
}
