package component

import "fyne.io/fyne/v2/widget"

func SelectorWidget(typesMap map[string]int, onSelect func(int)) *widget.Select {
	var typeNames []string
	for typeName := range typesMap {
		typeNames = append(typeNames, typeName)
	}

	typeSelect := widget.NewSelect(typeNames, func(selected string) {
		id := typesMap[selected]
		onSelect(id)
	})

	return typeSelect
}
