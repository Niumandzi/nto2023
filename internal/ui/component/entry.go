package component

import (
	"fyne.io/fyne/v2"
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
