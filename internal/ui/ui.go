package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	error2 "github.com/niumandzi/nto2023/internal/ui/page/error"
	"github.com/niumandzi/nto2023/internal/ui/page/event"
	"github.com/niumandzi/nto2023/internal/ui/page/index"
)

type GUI struct {
	App    fyne.App
	Window fyne.Window
	Event  event.EventPage
}

func NewGUI(app fyne.App, window fyne.Window, event event.EventPage) GUI {
	return GUI{
		App:    app,
		Window: window,
		Event:  event,
	}
}

func SetupUI(gui GUI) {
	w := gui.Window
	e := gui.Event

	mainContent := container.NewStack()

	mainContent.Add(index.ShowIndex())

	navBar := NavigationBar(e, mainContent, w)

	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}

func NavigationBar(event event.EventPage, mainContent *fyne.Container, window fyne.Window) *widget.Tree {
	treeData := map[string][]string{
		"": {"развлечения", "просвещение", "образование"}}

	navTree := widget.NewTreeWithStrings(treeData)
	navTree.OnSelected = func(id string) {
		var content fyne.CanvasObject

		switch id {
		case "развлечения":
			content = event.EventIndex("entertainment", window)
		case "просвещение":
			content = event.EventIndex("enlightenment", window)
		default:
			content = error2.ShowErrorPage()
		}

		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	return navTree
}
