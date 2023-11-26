package component

import "fyne.io/fyne/v2/widget"

func SelectorWidget(placeHolder string, typesMap interface{}, onSelectInt func(int), onSelectString func(string)) *widget.Select {
	var typeNames []string

	switch tMap := typesMap.(type) {
	case map[string]int:
		for typeName := range tMap {
			typeNames = append(typeNames, typeName)
		}
	case map[string]string:
		for typeName := range tMap {
			typeNames = append(typeNames, typeName)
		}
	default:
		// Handle error or unsupported types
	}

	typeSelect := widget.NewSelect(typeNames, func(selected string) {
		switch tMap := typesMap.(type) {
		case map[string]int:
			id := tMap[selected]
			if onSelectInt != nil {
				onSelectInt(id)
			}
		case map[string]string:
			value := tMap[selected]
			if onSelectString != nil {
				onSelectString(value)
			}
		}
	})
	typeSelect.PlaceHolder = placeHolder

	return typeSelect
}
