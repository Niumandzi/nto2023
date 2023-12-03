package work

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s WorkTypePage) CreateWorkType(window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWidget("Тип работ")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание нового типа работ", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateWorkType(nameEntry.Text, window, s.workTypeServ, onUpdate)
		}
	}, window)
}

func handleCreateWorkType(name string, window fyne.Window, workTypeServ service.WorkTypeService, onUpdate func()) {
	_, err := workTypeServ.CreateWorkType(name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип создан", "Тип работ успешно создан!", window)
		onUpdate()
	}
}
