package page

import "fyne.io/fyne/v2"

type EventPage interface {
	ShowEvent(category string) fyne.CanvasObject
	CreatEvent(category string) fyne.CanvasObject
	Update(category string) fyne.CanvasObject
}
