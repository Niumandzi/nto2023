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
	update := make(map[int]string)
	var delete []int
	var create []string

	vbox := container.NewVBox()
	nameEntry := component.EntryWidget("Помещение")
	nameEntry.SetText(name)
	vbox.Add(nameEntry)

	partsVBox := container.NewVBox()
	vbox.Add(partsVBox)

	for _, part := range parts {
		partID := part.ID
		partEntry := component.EntryWidget("Часть помещения")
		partEntry.SetText(part.Name)

		partEntry.OnChanged = func(newText string) {
			if newText != part.Name && newText != "" {
				update[partID] = newText
				println(1)
			} else {
				println(2)
				delete = append(delete, partID)
			}
		}

		partsVBox.Add(partEntry)
	}

	addPartButton := widget.NewButton("    Добавить часть для помещения    ", func() {
		newEntry := component.EntryWidget("Часть помещения")
		partsVBox.Add(newEntry)
		window.Canvas().Refresh(partsVBox)
	})

	vbox.Add(addPartButton)

	var customPopUp *widget.PopUp
	saveButton := widget.NewButton("            Обновить            ", func() {
		handleUpdateFacility(id, nameEntry.Text, update, delete, create, window, s.facilityServ, s.partServ, onUpdate, customPopUp)
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

func handleUpdateFacility(id int, name string, update map[int]string, delete []int, create []string, window fyne.Window, facilityServ service.FacilityService, partServ service.PartService, onUpdate func(), popUp *widget.PopUp) {
	err := facilityServ.UpdateFacility(id, name)

	if len(update) > 0 {
		err = partServ.UpdatePart(update)
	}

	if len(delete) > 0 {
		err = partServ.DeletePart(delete, false)
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
