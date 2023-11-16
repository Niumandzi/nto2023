package error

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowErrorPage() fyne.CanvasObject {
	label := widget.NewLabel("Этот раздел будет реализован позднее")
	return container.New(layout.NewCenterLayout(), label)
}
