package booking

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s BookingPage) CreateBooking(categoryName string, window fyne.Window, onUpdate func()) {
	partsEntries := make([]*widget.Entry, 0)
	vbox := container.NewVBox()

	var selectedEventID int
	var selectedFacilityID int

	events, err := s.eventServ.GetEvents(categoryName, 0, true)
	if err != nil {
		dialog.ShowError(err, window)
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	eventSelect := component.SelectorWidget("Мероприятие", eventNames, func(id int) {
		selectedEventID = id
	}, nil)
	vbox.Add(eventSelect)

	createDateLabel := widget.NewLabel(time.Now().Format("2006-02-01"))
	descriptionEntry := component.MultiLineEntryWidget("Описание")
	startDateEntry := component.EntryWidget("Дата начала (гггг-мм-дд)")
	endDateEntry := component.EntryWidget("Дата начала (гггг-мм-дд)")

	vbox.Add(createDateLabel)
	vbox.Add(descriptionEntry)
	vbox.Add(startDateEntry)
	vbox.Add(endDateEntry)

	facilities, err := s.facilityServ.GetFacilitiesByDate(startDateEntry.Text, endDateEntry.Text, true)
	if err != nil {
		dialog.ShowError(err, window)
	}

	facilityNames := make(map[string]int)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
	}

	facilitySelect := component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityID = id
	}, nil)

	vbox.Add(facilitySelect)

	partsVBox := container.NewVBox()

	addPartButton := widget.NewButton("    Добавить часть для помещения    ", func() {
		newEntry := component.EntryWidget("Часть помещения")
		partsEntries = append(partsEntries, newEntry)
		partsVBox.Add(newEntry)
		window.Canvas().Refresh(partsVBox)
	})

	deleteLastPartButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if len(partsEntries) > 0 {
			lastIndex := len(partsEntries) - 1
			partsVBox.Remove(partsEntries[lastIndex])
			partsEntries = partsEntries[:lastIndex]
			window.Canvas().Refresh(partsVBox)
		}
	})

	buttonBox := container.NewHBox(addPartButton, deleteLastPartButton)
	vbox.Add(buttonBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {
		formData := model.Booking{
			CreateDate:  time.Now().Format("2006-02-01"),
			Description: descriptionEntry.Text,
			StartDate:   startDateEntry.Text,
			EndDate:     endDateEntry.Text,
			EventID:     selectedEventID,
			FacilityID:  selectedFacilityID,
		}

		handleCreateBooking(formData, window, s.bookingServ, onUpdate, customPopUp)
	})

	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)
	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 100))
	customPopUp.Show()
}

func handleCreateBooking(appData model.Booking, window fyne.Window, bookingServ service.BookingService, onUpdate func(), popUp *widget.PopUp) {
	_, err := bookingServ.CreateBooking(appData)

	popUp.Hide()

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		dialog.ShowInformation("Бронирование создано", "Бронирование успешно создано!", window)
		onUpdate()
	}
}
