package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s FacilityPage) CreateFacility(window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWidget("Помещение")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание нового помещения", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateFacility(nameEntry.Text, window, s.facilityServ, onUpdate)
		}
	}, window)
}

func handleCreateFacility(name string, window fyne.Window, facilityServ service.FacilityService, onUpdate func()) {
	_, err := facilityServ.CreateFacility(name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Помещение создано", "Помещение успешно создано!", window)
		onUpdate()
	}
}
