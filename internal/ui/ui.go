package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/ui/page/application"
	"github.com/niumandzi/nto2023/internal/ui/page/booking"
	"github.com/niumandzi/nto2023/internal/ui/page/details"
	error2 "github.com/niumandzi/nto2023/internal/ui/page/error"
	"github.com/niumandzi/nto2023/internal/ui/page/event"
	"github.com/niumandzi/nto2023/internal/ui/page/facility"
	"github.com/niumandzi/nto2023/internal/ui/page/index"
	"github.com/niumandzi/nto2023/internal/ui/page/mug"
	"github.com/niumandzi/nto2023/internal/ui/page/registration"
	"github.com/niumandzi/nto2023/internal/ui/page/teacher"
	"github.com/niumandzi/nto2023/internal/ui/page/work"
)

type GUI struct {
	App    fyne.App
	Window fyne.Window
}

func NewGUI(app fyne.App, window fyne.Window) GUI {
	return GUI{
		App:    app,
		Window: window,
	}
}

func SetupUI(gui GUI, application application.ApplicationPage, booking booking.BookingPage, details details.DetailsPage, event event.EventPage, facility facility.FacilityPage, mugType mug.MugTypePage, teacher teacher.TeacherPage, registration registration.RegistrationPage, workType work.WorkTypePage) {
	w := gui.Window

	mainContent := container.NewStack()
	mainContent.Add(index.ShowIndex())

	navBar := NavigationBar(application, booking, details, event, facility, mugType, teacher, registration, workType, w, mainContent)

	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}

func NavigationBar(application application.ApplicationPage, booking booking.BookingPage, details details.DetailsPage, event event.EventPage, facility facility.FacilityPage, mugType mug.MugTypePage, teacher teacher.TeacherPage, registration registration.RegistrationPage, workType work.WorkTypePage, window fyne.Window, mainContent *fyne.Container) *widget.Tree {
	treeData := map[string][]string{
		"":             {"развлечения", "просвещение", "образование", "рабочий стол"},
		"развлечения":  {"типы развлечений", "работы развлечения", "бронь развлечения"},
		"просвещение":  {"типы просвещения", "работы просвещение", "бронь просвещение"},
		"образование":  {"типы кружка", "преподаватели", "регистрация кружков"},
		"рабочий стол": {"помещения", "типы работ"},
	}

	navTree := widget.NewTreeWithStrings(treeData)
	navTree.OnSelected = func(id string) {
		var content fyne.CanvasObject

		switch id {
		case "развлечения":
			content = event.IndexEvent("entertainment", window)
		case "просвещение":
			content = event.IndexEvent("enlightenment", window)
		case "типы развлечений":
			content = details.IndexDetails("entertainment", window)
		case "типы просвещения":
			content = details.IndexDetails("enlightenment", window)
		case "работы развлечения":
			content = application.IndexApplication("entertainment", "", window)
		case "работы просвещение":
			content = application.IndexApplication("enlightenment", "", window)
		case "бронь развлечения":
			content = booking.IndexBooking("entertainment", window)
		case "бронь просвещение":
			content = booking.IndexBooking("enlightenment", window)
		case "рабочий стол":
			content = application.IndexApplication("", "todo", window)
		case "помещения":
			content = facility.IndexFacility(window)
		case "типы работ":
			content = workType.IndexWorkType(window)
		case "образование":
			content = registration.ShowTableRegistration(window)
		case "преподаватели":
			content = teacher.IndexTeacher(window)
		case "типы кружка":
			content = mugType.IndexMugType(window)
		case "регистрация кружков":
			content = registration.IndexRegistration(window)
		default:
			content = error2.ShowErrorPage()
		}

		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	return navTree
}

//type EventPage interface {
//	IndexEvent(categoryName string, window fyne.Window) fyne.CanvasObject
//	ShowEvent(categoryName string, detailsID int, window fyne.Window, eventContainer *fyne.Container)
//	CreateEvent(categoryName string, window fyne.Window)
//	UpdateEvent(categoryName string, ID int, Name string, Date string, Description string, DetailsID int, window fyne.Window)
//}
