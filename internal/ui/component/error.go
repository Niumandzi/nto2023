package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ShowErrorDialog(err error, w fyne.Window) {
	dialog.ShowError(err, w)
}
