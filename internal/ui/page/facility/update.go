package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityPage) UpdateFacility(id int, name string, parts []model.Parts, window fyne.Window, onUpdate func()) {
	for _, part := range parts {
		println(part.ID, part.Name)
	}
	partsEntries := make([]*widget.Entry, 0)
	vbox := container.NewVBox()

	nameEntry := component.EntryWidget("Помещение")
	nameEntry.SetText(name)
	vbox.Add(nameEntry)

	partsVBox := container.NewVBox()
	vbox.Add(partsVBox)

	for _, part := range parts {
		partEntry := component.EntryWidget("Часть помещения")
		partEntry.SetText(part.Name)
		partsEntries = append(partsEntries, partEntry)
		partsVBox.Add(partEntry)
	}

	// Add/Delete part buttons similar to CreateFacility
	addPartButton := widget.NewButton("    Добавить часть для помещения    ", func() {
		// Similar implementation as in CreateFacility
	})
	deleteLastPartButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		// Similar implementation as in CreateFacility
	})
	buttonBox := container.NewHBox(addPartButton, deleteLastPartButton)
	vbox.Add(buttonBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Обновить            ", func() {
		var updatedParts []string
		for _, entry := range partsEntries {
			updatedParts = append(updatedParts, entry.Text)
		}
		//handleUpdateFacility(id, nameEntry.Text, updatedParts, window, s.facilityServ, onUpdate, customPopUp)
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

func handleUpdateFacility(id int, name string, parts []model.Parts, window fyne.Window, facilityServ service.FacilityService, onUpdate func()) {

	err := facilityServ.UpdateFacility(id, name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Помещение обновлено", "Помещение успешно обновлен!", window)
		onUpdate()
	}
}
