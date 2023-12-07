package teacher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (m TeacherPage) UpdateTeacher(id int, name string, window fyne.Window, onUpdate func()) {

	nameEntry := component.EntryWidgetWithData("Тип работы", name)

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Обновить информацию о преподаватель", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleUpdateEvent(id, nameEntry.Text, window, m.teacherServ, onUpdate)
		}
	}, window)
}

func handleUpdateEvent(id int, name string, window fyne.Window, teacherServ service.TeacherService, onUpdate func()) {

	err := teacherServ.UpdateTeacher(id, name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Преподаватель обновлен", "Тип кружка успешно обновлен!", window)
		onUpdate()
	}
}
