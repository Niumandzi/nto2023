package booking

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (b BookingPage) UpdateBooking(categoryName string, booking model.BookingWithFacility, window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()

	var selectedEventID int
	var selectedFacilityID int
	var facilityNames map[string]int
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select
	var partsBox *fyne.Container
	var customPopUp *widget.PopUp
	var selectedParts []int

	facilities, err := b.facilityServ.GetFacilitiesByDateTimeAndID(booking.StartDate, booking.StartTime, booking.EndDate, booking.EndTime, booking.ID, booking.Facility.ID)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	facilityNames = make(map[string]int)
	facilityParts = make(map[int]map[int]string)
	for _, facility := range facilities {
		facilityNames[facility.Name] = facility.ID
		partMap := make(map[int]string)
		for _, part := range facility.Parts {
			partMap[part.ID] = part.Name
		}
		facilityParts[facility.ID] = partMap
	}

	bookingPartIDs := make(map[int]bool)
	for _, part := range booking.Parts {
		bookingPartIDs[part.ID] = true
	}

	selectedParts = []int{}
	for partID := range bookingPartIDs {
		selectedParts = append(selectedParts, partID)
	}

	newPartsBox := container.NewVBox()
	for partID, partName := range facilityParts[booking.Facility.ID] {
		localPartID := partID

		isChecked := bookingPartIDs[localPartID]

		checkBox := widget.NewCheck(partName, nil)
		checkBox.Checked = isChecked
		checkBox.Refresh()

		checkBox.OnChanged = func(checked bool) {
			if checked {
				if !contains(selectedParts, localPartID) {
					selectedParts = append(selectedParts, localPartID)
					fmt.Print(selectedParts)
				}
			} else {
				for i, id := range selectedParts {
					if id == localPartID {
						selectedParts = append(selectedParts[:i], selectedParts[i+1:]...)
						break
					}
				}
			}
		}
		newPartsBox.Add(checkBox)
	}

	events, err := b.eventServ.GetActiveEvents(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	selectedEventID = booking.Event.ID
	eventSelect := component.SelectorWidget(booking.Event.Name, eventNames, func(id int) {
		selectedEventID = id
	}, nil)
	vbox.Add(eventSelect)

	createDateLabel := widget.NewLabel(booking.CreateDate)
	startLabel := widget.NewLabel("Дата/Время начала")
	endLabel := widget.NewLabel("Дата/Время окончания")
	descriptionEntry := component.MultiLineEntryWidgetWithData("Описание", booking.Description)
	startDateEntry := component.EntryWidgetWithData("гггг-мм-дд", booking.StartDate)
	startTimeEntry := component.EntryWidgetWithData("чч:мм", booking.StartTime)
	endDateEntry := component.EntryWidgetWithData("гггг-мм-дд", booking.EndDate)
	endTimeEntry := component.EntryWidgetWithData("чч:мм", booking.EndTime)

	saveButton := widget.NewButton("            Сохранить            ", func() {
		if facilityParts[selectedFacilityID] != nil && len(facilityParts[selectedFacilityID]) > 0 && len(selectedParts) == 0 {
			dialog.ShowError(fmt.Errorf("Для выбранного помещения необходимо выбрать хотя бы одну часть"), window)
			return
		}

		formData := model.Booking{
			ID:          booking.ID,
			CreateDate:  booking.CreateDate,
			Description: descriptionEntry.Text,
			StartDate:   startDateEntry.Text,
			StartTime:   startTimeEntry.Text,
			EndDate:     endDateEntry.Text,
			EndTime:     endTimeEntry.Text,
			EventID:     selectedEventID,
			FacilityID:  selectedFacilityID,
			PartIDs:     selectedParts,
		}

		handleUpdateBooking(formData, window, b.bookingServ, onUpdate, customPopUp)
	})
	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)

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

		vbox.Remove(buttons)
		vbox.Remove(partsBox)
		vbox.Add(newPartsBox)
		partsBox = newPartsBox
		vbox.Add(buttons)
		window.Canvas().Refresh(vbox)
	}

	vbox.Add(createDateLabel)
	vbox.Add(descriptionEntry)
	startContainer := container.New(layout.NewGridLayoutWithColumns(2), startDateEntry, startTimeEntry)
	vbox.Add(startLabel)
	vbox.Add(startContainer)
	endContainer := container.New(layout.NewGridLayoutWithColumns(2), endDateEntry, endTimeEntry)
	vbox.Add(endLabel)
	vbox.Add(endContainer)

	facilityNames = make(map[string]int)
	updateFacilities := func() {
		if validateDate(startDateEntry.Text) && validateTime(startTimeEntry.Text) && validateDate(startDateEntry.Text) && validateTime(endTimeEntry.Text) {
			facilities, err = b.facilityServ.GetFacilitiesByDateTimeAndID(booking.StartDate, booking.StartTime, booking.EndDate, booking.EndTime, booking.ID, booking.Facility.ID)
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

			for _, facility := range facilities {
				parts := make(map[int]string)
				for _, part := range facility.Parts {
					parts[part.ID] = part.Name
				}
				facilityParts[facility.ID] = parts
			}

			if facilitySelect != nil {
				vbox.Remove(partsBox)
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

	updateFacilities()

	startDateEntry.OnChanged = func(string) { updateFacilities() }
	startTimeEntry.OnChanged = func(string) { updateFacilities() }
	endDateEntry.OnChanged = func(string) { updateFacilities() }
	endTimeEntry.OnChanged = func(string) { updateFacilities() }

	selectedFacilityID = booking.Facility.ID
	facilitySelect = component.SelectorWidget(booking.Facility.Name, facilityNames, func(id int) {
		selectedFacilityID = id
		updateParts()
	}, nil)

	vbox.Add(facilitySelect)
	partsBox = newPartsBox
	vbox.Add(partsBox)

	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 300))
	customPopUp.Show()
}

func handleUpdateBooking(formDate model.Booking, window fyne.Window, bookingServ service.BookingService, onUpdate func(), popUp *widget.PopUp) {
	err := bookingServ.UpdateBooking(formDate)

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		popUp.Hide()
		dialog.ShowInformation("Бронирование сохранено", "Бронирование успешно сохранено!", window)
		onUpdate()
	}
}
