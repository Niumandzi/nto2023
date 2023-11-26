package work

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s WorkTypePage) UpdateWorkType(id int, name string, window fyne.Window, onUpdate func()) {

	nameEntry := component.EntryWithDataWidget("Тип работы", name)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Обновить тип работ", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, nameEntry.Text, window, s.workTypeServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(id int, name string, window fyne.Window, workTypeServ service.WorkTypeService, onUpdate func()) {

	err := workTypeServ.UpdateWorkType(id, name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип обновлен", "Тип работ успешно обновлен!", window)
		onUpdate()
	}
}
