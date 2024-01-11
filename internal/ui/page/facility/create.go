package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s FacilityPage) CreateFacility(window fyne.Window, onUpdate func()) {
	partsEntries := make([]*widget.Entry, 0)
	vbox := container.NewVBox()

	nameEntry := component.EntryWidget("Помещение")
	vbox.Add(nameEntry)

	partsVBox := container.NewVBox()
	vbox.Add(partsVBox)

	addPartButton := widget.NewButton("     Добавить части для помещения    ", func() {
		for i := 0; i < 2; i++ {
			if len(partsEntries) >= 2 {
				return
			}
			newEntry := component.EntryWidget("Часть помещения")
			partsEntries = append(partsEntries, newEntry)
			partsVBox.Add(newEntry)
			window.Canvas().Refresh(partsVBox)
		}
	})

	deleteLastPartButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		for i := 0; i < 2; i++ {
			if len(partsEntries) > 0 {
				lastIndex := len(partsEntries) - 1
				partsVBox.Remove(partsEntries[lastIndex])
				partsEntries = partsEntries[:lastIndex]
				window.Canvas().Refresh(partsVBox)
			}
		}
	})

	buttonBox := container.NewHBox(addPartButton, deleteLastPartButton)
	vbox.Add(buttonBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {
		var parts []string
		for _, entry := range partsEntries {
			parts = append(parts, entry.Text)
		}
		handleCreateFacility(nameEntry.Text, parts, window, s.facilityServ, onUpdate, customPopUp)
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

func handleCreateFacility(name string, parts []string, window fyne.Window, facilityServ service.FacilityService, onUpdate func(), popUp *widget.PopUp) {
	_, err := facilityServ.CreateFacility(name, parts)

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		popUp.Hide()
		dialog.ShowInformation("Помещение создано", "Помещение успешно создано!", window)
		onUpdate()
	}
}
