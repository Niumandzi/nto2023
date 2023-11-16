package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/niumandzi/nto2023/internal/ui/component"
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

	// Создание основного контента
	mainContent := container.NewMax()

	// Создание панели навигации
	navBar := component.NavigationBar(mainContent, w) // Предположим, что 'cases' доступны

	// Размещение панели навигации и основного контента
	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2 // Установка ширины панели навигации

	// Установка контента окна
	w.SetContent(split)

	// Отображение окна
	w.ShowAndRun()
}
