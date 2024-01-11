package mug

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (m MugTypePage) UpdateMugType(id int, name string, window fyne.Window, onUpdate func()) {

	nameEntry := component.EntryWidgetWithData("Тип работы", name)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Обновить тип работ", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, nameEntry.Text, window, m.mugTypeServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(id int, name string, window fyne.Window, mugTypeServ service.MugTypeService, onUpdate func()) {

	err := mugTypeServ.UpdateMugType(id, name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип обновлен", "Тип кружка успешно обновлен!", window)
		onUpdate()
	}
}
