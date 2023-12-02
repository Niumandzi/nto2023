package booking

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (s BookingPage) CreateBooking(categoryName string, window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()

	var selectedEventID int
	var selectedFacilityID int
	var facilityNames map[string]int
	var facilitySelect *widget.Select

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
	endDateEntry := component.EntryWidget("Дата окончания (гггг-мм-дд)")

	vbox.Add(createDateLabel)
	vbox.Add(descriptionEntry)
	vbox.Add(startDateEntry)
	vbox.Add(endDateEntry)

	updateFacilities := func() {
		if startDateEntry.Text != "" && endDateEntry.Text != "" {
			facilities, err := s.facilityServ.GetFacilitiesByDate(startDateEntry.Text, endDateEntry.Text)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			facilityNames = make(map[string]int)
			for _, facility := range facilities {
				facilityNames[facility.Name] = facility.ID
				partsDescription := ""
				for _, part := range facility.Parts {
					if partsDescription != "" {
						partsDescription += ", "
					}

					partsDescription += part.Name
				}

				if facility.HaveParts && partsDescription != "" {
					facilityNames[facility.Name+" Части: "+partsDescription] = facility.ID
				} else {
					facilityNames[facility.Name] = facility.ID
				}
			}

			if facilitySelect != nil {
				facilitySelect.Options = getFacilityOptions(facilityNames)
				facilitySelect.Refresh()
			}
		}

	}

	startDateEntry.OnChanged = func(string) { updateFacilities() }
	endDateEntry.OnChanged = func(string) { updateFacilities() }

	facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityID = id
	}, nil)

	vbox.Add(facilitySelect)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {
		var parts []int

		formData := model.Booking{
			CreateDate:  time.Now().Format("2006-02-01"),
			Description: descriptionEntry.Text,
			StartDate:   startDateEntry.Text,
			EndDate:     endDateEntry.Text,
			EventID:     selectedEventID,
			FacilityID:  selectedFacilityID,
			PartIDs:     parts,
		}

		handleCreateBooking(formData, window, s.bookingServ, onUpdate, customPopUp)
	})

	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)
	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 300))
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

func getFacilityOptions(facilityNames map[string]int) []string {
	var options []string
	for name := range facilityNames {
		options = append(options, name)
	}
	return options
}
