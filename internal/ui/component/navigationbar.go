package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/ui/page/entertainment"
)

func NavigationBar(mainContent *fyne.Container, window fyne.Window) *widget.Tree {
	treeData := map[string][]string{
		"": {"развлечение", "просвещение", "образование"}}

	navTree := widget.NewTreeWithStrings(treeData)
	navTree.OnSelected = func(id string) {
		var content fyne.CanvasObject
		// Обработка выбранного элемента
		switch id {
		case "развлечение":
			content = entertainment.ShowIntertainment()
		default:
			content = widget.NewLabel("Выберите категорию")
		}

		mainContent.Objects = []fyne.CanvasObject{content}
		mainContent.Refresh() // Обновляем содержимое контейнера
	}

	return navTree
}
