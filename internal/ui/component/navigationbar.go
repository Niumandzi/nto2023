package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	error2 "github.com/niumandzi/nto2023/internal/ui/page/error"
	"github.com/niumandzi/nto2023/internal/ui/page/event"
)

func NavigationBar(mainContent *fyne.Container, window fyne.Window) *widget.Tree {
	treeData := map[string][]string{
		"": {"развлечения", "просвещение", "образование"}}

	navTree := widget.NewTreeWithStrings(treeData)
	navTree.OnSelected = func(id string) {
		var content fyne.CanvasObject

		switch id {
		case "развлечения":
			content = event.ShowEvent("entertainment")
		case "просвещение":
			content = event.ShowEvent("enlightenment")
		default:
			content = error2.ShowErrorPage()
		}

		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh()
	}

	return navTree
}
