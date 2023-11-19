package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (s DetailsPage) UpdateDetail(id int, category string, typeName string, window fyne.Window) {
	switch category {
	case "entertainment":
		category = "Развлечения"
	case "enlightenment":
		category = "Просвещение"
	case "education":
		category = "Образование"
	}

	categoryLabel := widget.NewLabel(category)
	typeNameEntry := component.EntryWithDataWidget("", typeName)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", categoryLabel),
		widget.NewFormItem("", typeNameEntry),
	}

	dialog.ShowForm("Обновить событие", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, typeNameEntry.Text, window, s.detailsServ)
		}
	}, window)
}

func handleUpdateEvent(detailID int, typeName string, window fyne.Window, detailsServ service.DetailsService) {

	err := detailsServ.UpdateDetail(detailID, typeName)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Событие обновлено", "Событие успешно обновлено!", window)
	}
}