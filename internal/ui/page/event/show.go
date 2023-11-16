package event

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	error2 "github.com/niumandzi/nto2023/internal/ui/page/error"
)

func ShowEvent(category string) fyne.CanvasObject {
	var content fyne.CanvasObject

	switch category {
	case "entertainment":
		content = widget.NewLabel("Добро пожаловать в раздел развлечения")
	case "enlightenment":
		content = widget.NewLabel("Добро пожаловать в раздел просвещение")
	default:
		content = error2.ShowErrorPage()
	}

	return content
}
