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
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select

	events, err := s.eventServ.GetActiveEvents(categoryName)
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

	var selectedParts []int

	var partsBox *fyne.Container
	partsBox = container.NewVBox()
	vbox.Add(partsBox)

	updateParts := func() {
		newPartsBox := container.NewVBox()
		for partID, partName := range facilityParts[selectedFacilityID] {
			localPartID := partID
			checkBox := widget.NewCheck(partName, func(checked bool) {
				if checked {
					if !contains(selectedParts, localPartID) {
						selectedParts = append(selectedParts, localPartID)
					}
				} else {
					for i, id := range selectedParts {
						if id == localPartID {
							selectedParts = append(selectedParts[:i], selectedParts[i+1:]...)
							break
						}
					}
				}
			})
			newPartsBox.Add(checkBox)
		}

		vbox.Remove(partsBox)
		vbox.Add(newPartsBox)
		partsBox = newPartsBox
		window.Canvas().Refresh(vbox)
	}

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {

		formData := model.Booking{
			CreateDate:  time.Now().Format("2006-02-01"),
			Description: descriptionEntry.Text,
			StartDate:   startDateEntry.Text,
			EndDate:     endDateEntry.Text,
			EventID:     selectedEventID,
			FacilityID:  selectedFacilityID,
			PartIDs:     selectedParts,
		}

		handleCreateBooking(formData, window, s.bookingServ, onUpdate, customPopUp)
	})

	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)

	vbox.Add(createDateLabel)
	vbox.Add(descriptionEntry)
	vbox.Add(startDateEntry)
	vbox.Add(endDateEntry)

	facilityNames = make(map[string]int)

	updateFacilities := func() {
		if startDateEntry.Text != "" && endDateEntry.Text != "" {
			facilities, err := s.facilityServ.GetFacilitiesByDate(startDateEntry.Text, endDateEntry.Text)
			if err != nil {
				dialog.ShowError(err, window)
			}

			for key, _ := range facilityNames {
				delete(facilityNames, key)
			}

			for _, facility := range facilities {
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

			facilityParts = make(map[int]map[int]string)
			for _, facility := range facilities {
				parts := make(map[int]string)
				for _, part := range facility.Parts {
					parts[part.ID] = part.Name
				}
				facilityParts[facility.ID] = parts
			}

			if facilitySelect != nil {
				vbox.Remove(facilitySelect)
				vbox.Remove(buttons)
				facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
					selectedFacilityID = id
					updateParts()
				}, nil)
				vbox.Add(facilitySelect)
				vbox.Add(buttons)
			}

		}
	}

	startDateEntry.OnChanged = func(string) { updateFacilities() }
	endDateEntry.OnChanged = func(string) { updateFacilities() }

	facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityID = id
		updateParts()
	}, nil)

	vbox.Add(facilitySelect)

	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 300))
	customPopUp.Show()
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
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
