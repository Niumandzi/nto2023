package details

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s DetailsPage) UpdateDetail(id int, category string, typeName string, window fyne.Window, onUpdate func()) {
	switch category {
	case "entertainment":
		category = "Развлечения"
	case "enlightenment":
		category = "Просвещение"
	case "education":
		category = "Образование"
	}

	categoryLabel := widget.NewLabel(category)
	typeNameEntry := component.EntryWidgetWithData("Тип события", typeName)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", categoryLabel),
		widget.NewFormItem("", typeNameEntry),
	}

	dialog.ShowForm("Обновить событие", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, typeNameEntry.Text, window, s.detailsServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(detailID int, typeName string, window fyne.Window, detailsServ service.DetailsService, onUpdate func()) {

	err := detailsServ.UpdateDetail(detailID, typeName)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Тип обновлен", "Тип успешно обновлен!", window)
		onUpdate()
	}
}
