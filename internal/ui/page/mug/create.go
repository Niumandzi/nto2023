package mug

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (m MugTypePage) CreateMugType(window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWidget("Тип кружка")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание нового типа кружка", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateMugType(nameEntry.Text, window, m.mugTypeServ, onUpdate)
		}
	}, window)
}

func handleCreateMugType(name string, window fyne.Window, mugTypeServ service.MugTypeService, onUpdate func()) {
	_, err := mugTypeServ.CreateMugType(name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип создан", "Тип кружка успешно создан!", window)
		onUpdate()
	}
}
