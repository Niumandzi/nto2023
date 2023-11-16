package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ShowErrorDialogWidget(err error, w fyne.Window) {
	dialog.ShowError(err, w)
}
