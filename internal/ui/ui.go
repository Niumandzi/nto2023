package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/internal/ui/page/index"
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

func SetupUI(gui GUI) {
	w := gui.Window

	mainContent := container.NewStack()

	mainContent.Add(index.ShowIndex())

	navBar := component.NavigationBar(mainContent, w)

	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}
