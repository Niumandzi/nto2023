package component

import "fyne.io/fyne/v2/widget"

func TypeSelectWidget(types map[string]string, onSelected func(string)) *widget.Select {
	// Получаем список русских терминов для отображения в выпадающем списке
	typeNames := make([]string, 0, len(types))
	for name := range types {
		typeNames = append(typeNames, name)
	}

	// Создаем выпадающее меню с русскими терминами
	typeSelect := widget.NewSelect(typeNames, func(selected string) {
		// При выборе вызываем callback с соответствующим английским ключом
		onSelected(types[selected])
	})

	return typeSelect
}
