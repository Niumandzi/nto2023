package teacher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
)

func (t TeacherPage) CreateTeacher(window fyne.Window, onUpdate func()) {
	nameEntry := component.EntryWidget("Преподаватель")

	formItems := []*widget.FormItem{
		widget.NewFormItem("", nameEntry),
	}

	dialog.ShowForm("Создание записи о преподавателе", "Создать", "Отмена", formItems, func(confirm bool) {
		if confirm {
			handleCreateTeacher(nameEntry.Text, window, t.teacherServ, onUpdate)
		}
	}, window)
}

func handleCreateTeacher(name string, window fyne.Window, teacherServ service.TeacherService, onUpdate func()) {
	_, err := teacherServ.CreateTeacher(name)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Запись создан", "Запись успешно создана!", window)
		onUpdate()
	}
}
