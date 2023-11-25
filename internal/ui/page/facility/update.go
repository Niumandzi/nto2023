package facility

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s FacilityPage) UpdateFacility(id int, name string, window fyne.Window, onUpdate func()) {

	nameEntry := component.EntryWithDataWidget("Помещение", name)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Обновить помещение", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, nameEntry.Text, window, s.facilityServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(id int, name string, window fyne.Window, facilityServ service.FacilityService, onUpdate func()) {

	err := facilityServ.UpdateFacility(id, name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Помещение обновлено", "Помещение успешно обновлен!", window)
		onUpdate()
	}
}
