package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/ui/page/details"
	error2 "github.com/niumandzi/nto2023/internal/ui/page/error"
	"github.com/niumandzi/nto2023/internal/ui/page/event"
	"github.com/niumandzi/nto2023/internal/ui/page/facility"
	"github.com/niumandzi/nto2023/internal/ui/page/index"
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

func SetupUI(gui GUI, event event.EventPage, details details.DetailsPage, facility facility.FacilityPage, workType work.WorkTypePage) {
	w := gui.Window

	mainContent := container.NewStack()

	mainContent.Add(index.ShowIndex())

	navBar := NavigationBar(event, details, facility, workType, mainContent, w)

	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}

func NavigationBar(event event.EventPage, details details.DetailsPage, facility facility.FacilityPage, workType work.WorkTypePage, mainContent *fyne.Container, window fyne.Window) *widget.Tree {
	treeData := map[string][]string{
		"":             {"развлечения", "просвещение", "образование", "рабочий стол"},
		"развлечения":  {"типы развлечений"},
		"просвещение":  {"типы просвещения"},
		"образование":  {"типы образования"},
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
		case "помещения":
			content = facility.IndexFacility(window)
		case "типы работ":
			content = workType.IndexWorkType(window)
		default:
			content = error2.ShowErrorPage()
		}

		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	return navTree
}

type EventPage interface {
	IndexEvent(categoryName string, window fyne.Window) fyne.CanvasObject
	ShowEvent(categoryName string, detailsID int, window fyne.Window, eventContainer *fyne.Container)
	CreateEvent(categoryName string, window fyne.Window)
	UpdateEvent(categoryName string, ID int, Name string, Date string, Description string, DetailsID int, window fyne.Window)
}

type DetailsPage interface {
	IndexDetails(categoryName string, window fyne.Window) fyne.CanvasObject
	ShowDetails(categoryName string, window fyne.Window, eventContainer *fyne.Container)
}
