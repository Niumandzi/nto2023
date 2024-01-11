package teacher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type TeacherPage struct {
	teacherServ service.TeacherService
	logger      logging.Logger
}

func NewTeacherPage(teach service.TeacherService, logger logging.Logger) TeacherPage {
	return TeacherPage{
		teacherServ: teach,
		logger:      logger,
	}
}

func (t TeacherPage) IndexTeacher(window fyne.Window) fyne.CanvasObject {
	teacherContainer := container.NewStack()
	teacherList := func(eventType string) {
		t.ShowTeacher(window, teacherContainer)
	}

	createTeacherButton := widget.NewButton("Создать новую запись о преподавателе", func() {
		t.CreateTeacher(window, func() {
			teacherList("")
		})
	})

	createButtons := container.NewHBox(createTeacherButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, teacherContainer)
	teacherList("")

	return content
}
