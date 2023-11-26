package details

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s DetailsPage) CreateDetails(categoryName string, window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWidget("Тип события")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание нового типа события", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateDetails(nameEntry.Text, categoryName, window, s.detailsServ, onUpdate)
		}
	}, window)
}

func handleCreateDetails(eventName string, categoryName string, window fyne.Window, detailsServ service.DetailsService, onUpdate func()) {
	_, err := detailsServ.CreateDetail(categoryName, eventName)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип создан", "Тип для события успешно создано!", window)
		onUpdate()
	}
}
