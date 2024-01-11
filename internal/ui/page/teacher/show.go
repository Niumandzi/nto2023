package teacher

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
)

func (t TeacherPage) ShowTeacher(window fyne.Window, eventContainer *fyne.Container) {
	teachers, err := t.teacherServ.GetTeachers()
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	eventContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, teacher := range teachers {
		card := t.createTeacherCard(teacher, window, func() {
			t.ShowTeacher(window, eventContainer)
		})
		grid.Add(card)
	}

	eventContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	eventContainer.Refresh()
}

func (t TeacherPage) createTeacherCard(teacher model.Teacher, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText, isActive := card(teacher)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		t.UpdateTeacher(teacher.ID, teacher.Name, window, onUpdate)
	})

	var icon fyne.Resource
	var dialogTitle, dialogMessage string

	if isActive {
		icon = theme.DeleteIcon()
		dialogTitle = "Преподаватель удален"
		dialogMessage = "Преподаватель успешно удален!"
	} else {
		icon = theme.ViewRefreshIcon()
		dialogTitle = "Преподаватель восстановлен"
		dialogMessage = "Преподаватель успешно восстановлен!"
	}

	deleteButton := widget.NewButtonWithIcon("", icon, func() {
		err := t.teacherServ.DeleteRestoreTeacher(teacher.ID, !isActive)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation(dialogTitle, dialogMessage, window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)
	eventContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return eventContainer
}

func card(teacher model.Teacher) (string, bool) {
	return fmt.Sprintf("Имя: %s", teacher.Name), teacher.IsActive
}
