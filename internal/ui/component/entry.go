package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func EntryWidget(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	return entry
}

func MultiLineEntryWidget(placeholder string) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder(placeholder)
	entry.Wrapping = fyne.TextWrapWord
	return entry
}

func EntryWithDataWidget(placeholder string, data string) *widget.Entry {
	dataBinding := binding.NewString()
	_ = dataBinding.Set(data)

	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)

	initialData, err := dataBinding.Get()
	if err == nil {
		entry.SetText(initialData)
	}

	entry.OnChanged = func(newText string) {
		_ = dataBinding.Set(newText)
	}

	return entry
}

func MultiLineEntryWidgetWithData(placeholder string, data string) *widget.Entry {
	dataBinding := binding.NewString()
	_ = dataBinding.Set(data)

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder(placeholder)
	entry.Wrapping = fyne.TextWrapWord

	initialData, err := dataBinding.Get()
	if err == nil {
		entry.SetText(initialData)
	}

	entry.OnChanged = func(newText string) {
		_ = dataBinding.Set(newText)
	}

	return entry
}
