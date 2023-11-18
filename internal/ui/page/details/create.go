package event

//func (s DetailsPage) CreateEventType(categoryName string, window fyne.Window) {
//	nameEntry := component.EntryWidget("Тип события")
//
//	formItems := []*widget.FormItem{
//		widget.NewFormItem("", nameEntry),
//	}
//
//	dialog.ShowForm("Создание нового типа события", "Создать", "Отмена", formItems, func(confirm bool) {
//		if confirm {
//			handleCreateDetails(nameEntry.Text, categoryName, window, s.eventServ)
//		}
//	}, window)
//}
//
//func handleCreateDetails(eventName string, categoryName string, window fyne.Window, eventServ service.DetailsServ) {
//	_, err := eventServ.CreateDetails(categoryName, eventName)
//	if err != nil {
//		dialog.ShowError(err, window)
//	} else {
//		dialog.ShowInformation("Тип создан", "Тип для события успешно создано!", window)
//	}
//}
