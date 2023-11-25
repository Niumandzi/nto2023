package application

//
//import (
//	"fyne.io/fyne/v2"
//	"fyne.io/fyne/v2/dialog"
//	"fyne.io/fyne/v2/widget"
//	"github.com/niumandzi/nto2023/internal/service"
//	"github.com/niumandzi/nto2023/internal/ui/component"
//	"github.com/niumandzi/nto2023/model"
//	"time"
//)
//
//func (s ApplicationPage) CreateApplication1(categoryName string, window fyne.Window, onUpdate func()) {
//	formData := struct {
//		Name        string
//		Date        string
//		Description string
//		DetailsID   int
//	}{}
//
//	nameEntry := component.EntryWidget("Название")
//	dateEntry := component.EntryWidget("дд.мм.гггг")
//	descriptionEntry := component.MultiLineEntryWidget("Описание")
//
//	details, err := s.applicationServ.GetDetails(categoryName)
//	if err != nil {
//		dialog.ShowError(err, window)
//	}
//
//	typeNames := make(map[string]int)
//	for _, detail := range details {
//		typeNames[detail.TypeName] = detail.ID
//	}
//
//	detailsSelect := component.SelectorWidget("Тип", typeNames, func(id int) {
//		formData.DetailsID = id
//	})
//
//	formItems := []*widget.FormItem{
//		widget.NewFormItem("", detailsSelect),
//		widget.NewFormItem("", nameEntry),
//		widget.NewFormItem("", dateEntry),
//		widget.NewFormItem("", descriptionEntry),
//	}
//
//	dialog.ShowForm("                            Создать событие                           ", "Создать", "Отмена", formItems, func(confirm bool) {
//
//		formData.Name = nameEntry.Text
//		formData.Date = dateEntry.Text
//		formData.Description = descriptionEntry.Text
//
//		if confirm {
//			handleCreateApplication(formData.Name, formData.Date, formData.Description, formData.DetailsID, window, s.applicationServ, onUpdate)
//		}
//	}, window)
//}
//
//func (s ApplicationPage) CreateApplication(categoryName string, window fyne.Window, onUpdate func()) {
//	formData := struct {
//		Description string
//		CreateDate  string
//		DueDate     string
//		Status      string
//		WorkTypeId  int
//		FacilityId int
//		EventId    int
//	}{}
//
//	createDate := time.Now().Format("02-01-2006") // Getting current date in dd-mm-yyyy format
//	descriptionEntry := component.MultiLineEntryWidget("Описание")
//	dueDateEntry := component.EntryWidget("Срок выполнения (дд.мм.гггг)")
//	statusEntry := component.EntryWidget("Статус")
//
//	facilities, err := s.facilityServ.GetFacilities()
//	if err != nil {
//		dialog.ShowError(err, window)
//	}
//
//	facilityNames := make(map[string]int)
//	for _, facility := range facilities {
//		facilityNames[facility.Name] = facility.ID
//	}
//
//	facilitySelect := component.SelectorWidget("Помещение", facilityNames, func(id int) {
//		formData.FacilityId
//	})
//
//	workTypes, err := s.workTypeServ.GetWorkTypes()
//	if err != nil {
//		dialog.ShowError(err, window)
//	}
//
//	workNames := make(map[string]int)
//	for _, work := range workTypes {
//		workNames[work.Name] = work.ID
//	}
//
//	workSelect := component.SelectorWidget("Тип работ", workNames, updateApplicationList)
//
//	formItems := []*widget.FormItem{
//		widget.NewFormItem("Описание", descriptionEntry),
//		widget.NewFormItem("Срок выполнения", dueDateEntry),
//		widget.NewFormItem("Статус", statusEntry),
//		// Add form items for WorkTypeId, EventId, FacilityId here
//	}
//
//	dialog.ShowForm("                            Создать событие                           ", "Создать", "Отмена", formItems, func(confirm bool) {
//		if confirm {
//			formData.Description = descriptionEntry.Text
//			formData.DueDate = dueDateEntry.Text
//			formData.Status = statusEntry.Text
//			// Set values for WorkTypeId, EventId, FacilityId
//			// ...
//
//			handleCreateApplication(createDate, formData, window, s.applicationServ, onUpdate)
//		}
//	}, window)
//}
//
//func handleCreateApplication(createDate string, formData struct {
//	Description string
//	CreateDate  string
//	DueDate     string
//	Status      string
//	WorkTypeId  int
//	EventId     int
//	FacilityId  int
//}, window fyne.Window, applicationServ service.ApplicationService, onUpdate func()) {
//	newApplication := model.Application{
//		CreateDate:  createDate,
//		Description: formData.Description,
//		DueDate:     formData.DueDate,
//		Status:      formData.Status,
//		WorkTypeId:  formData.WorkTypeId,
//		EventId:     formData.EventId,
//		FacilityId:  formData.FacilityId,
//	}
//
//	_, err := applicationServ.CreateApplication(newApplication)
//	if err != nil {
//		dialog.ShowError(err, window)
//	} else {
//		dialog.ShowInformation("Событие создано", "Событие успешно создано!", window)
//		onUpdate()
//	}
//}
