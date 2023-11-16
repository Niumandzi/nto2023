package index

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowIndex() fyne.CanvasObject {
	label := widget.NewLabel("Your Text Here")
	return container.New(layout.NewCenterLayout(), label)
}
