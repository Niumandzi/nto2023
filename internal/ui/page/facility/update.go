package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
	"image/color"
)

func (s FacilityPage) UpdateFacility(id int, name string, parts []model.Part, window fyne.Window, onUpdate func()) {
	partsEntries := make([]*widget.Entry, 0)
	updateData := make(map[int]string)
	deleteData := make(map[int]bool)

	vbox := container.NewVBox()
	nameEntry := component.EntryWidgetWithData("Помещение", name)
	vbox.Add(nameEntry)

	partsVBox := container.NewVBox()
	partsVBox2 := container.NewVBox()
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(0, 36))
	partsVBox2.Add(spacer)
	vbox.Add(partsVBox)

	for _, part := range parts {
		partID := part.ID
		isActive := part.IsActive

		partEntry := component.EntryWidgetWithData("Часть помещения", part.Name)
		partEntry.Resize(fyne.NewSize(100, 36))

		deleteButton := widget.NewButtonWithIcon("", nil, nil)
		if isActive {
			deleteButton.SetIcon(theme.DeleteIcon())
		} else {
			deleteButton.SetIcon(theme.MediaReplayIcon())
		}

		deleteButton.OnTapped = func() {
			isActive = !isActive
			deleteData[partID] = isActive

			if isActive {
				deleteButton.SetIcon(theme.DeleteIcon())
			} else {
				deleteButton.SetIcon(theme.MediaReplayIcon())
			}

			window.Canvas().Refresh(partsVBox)
		}

		partEntry.OnChanged = func(newText string) {
			if newText != part.Name && newText != "" {
				updateData[partID] = newText
			}
		}

		partsVBox.Add(partEntry)
		partsVBox2.Add(deleteButton)
	}

	addPartButton := widget.NewButton("      Добавить часть для помещения      ", func() {
		newEntry := component.EntryWidget("Часть помещения")
		partsEntries = append(partsEntries, newEntry)
		partsVBox.Add(newEntry)
		window.Canvas().Refresh(partsVBox)
	})

	deleteLastPartButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if len(partsEntries) > 0 {
			lastIndex := len(partsEntries) - 1
			partsVBox.Remove(partsEntries[lastIndex])
			partsEntries = partsEntries[:lastIndex]
			window.Canvas().Refresh(partsVBox)
		}
	})

	buttonBox := container.NewHBox(addPartButton, deleteLastPartButton)
	vbox.Add(buttonBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Обновить            ", func() {
		var createData []string
		for _, entry := range partsEntries {
			createData = append(createData, entry.Text)
		}
		handleUpdateFacility(id, nameEntry.Text, updateData, deleteData, createData, window, s.facilityServ, s.partServ, onUpdate, customPopUp)
	})

	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)
	vbox.Add(buttons)
	testBox := container.NewHBox(vbox, partsVBox2)

	customPopUp = widget.NewModalPopUp(testBox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 100))
	customPopUp.Show()
}

func handleUpdateFacility(id int, name string, update map[int]string, delete map[int]bool, create []string, window fyne.Window, facilityServ service.FacilityService, partServ service.PartService, onUpdate func(), popUp *widget.PopUp) {
	err := facilityServ.UpdateFacility(id, name)

	if len(update) > 0 {
		err = partServ.UpdatePart(update)
	}

	if len(delete) > 0 {
		err = partServ.DeletePart(delete)
	}

	if len(create) > 0 {
		_, err = partServ.CreatePart(id, create)
	}

	popUp.Hide()

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Помещение обновлено", "Помещение успешно обновлено!", window)
		onUpdate()
	}
}
