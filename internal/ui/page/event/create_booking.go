package event

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (e EventPage) CreateBooking(selectedEventID int, eventName string, categoryName string, window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()

	var selectedFacilityID int
	var facilityNames map[string]int
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select

	events, err := e.eventServ.GetActiveEvents(categoryName)
	if err != nil {
		dialog.ShowError(err, window)
	}

	eventNames := make(map[string]int)
	for _, event := range events {
		eventNames[event.Name] = event.ID
	}

	eventSelect := widget.NewLabel(eventName)
	vbox.Add(eventSelect)

	createDateLabel := widget.NewLabel(time.Now().Format("2006-02-01"))
	descriptionEntry := component.MultiLineEntryWidget("Описание")
	startDateEntry := component.EntryWidget("Дата начала (гггг-мм-дд)")
	startTimeEntry := component.EntryWidget("Время начала (чч:мм)")
	endDateEntry := component.EntryWidget("Дата окончания (гггг-мм-дд)")
	endTimeEntry := component.EntryWidget("Время начала (чч:мм)")

	var selectedParts []int

	var partsBox *fyne.Container
	partsBox = container.NewVBox()
	vbox.Add(partsBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {
		if facilityParts[selectedFacilityID] != nil && len(facilityParts[selectedFacilityID]) > 0 && len(selectedParts) == 0 {
			dialog.ShowError(fmt.Errorf("для выбранного помещения необходимо выбрать хотя бы одну часть"), window)
		}

		formData := model.Booking{
			CreateDate:  time.Now().Format("2006-02-01"),
			Description: descriptionEntry.Text,
			StartDate:   startDateEntry.Text,
			StartTime:   startTimeEntry.Text,
			EndDate:     endDateEntry.Text,
			EndTime:     endTimeEntry.Text,
			EventID:     selectedEventID,
			FacilityID:  selectedFacilityID,
			PartIDs:     selectedParts,
		}

		if len(selectedParts) < len(facilityParts[selectedFacilityID]) {
			infoDialog := dialog.NewCustom("Частичное бронирование", "OK", widget.NewLabel("Зал будет забронирован частично, так как не все доступные части выбраны."), window)
			infoDialog.SetOnClosed(func() {
				handleCreateBooking(formData, window, e.bookingServ, onUpdate, customPopUp)
			})
			infoDialog.Show()
		}
		handleCreateBooking(formData, window, e.bookingServ, onUpdate, customPopUp)
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
	vbox.Add(startDateEntry)
	vbox.Add(startTimeEntry)
	vbox.Add(endDateEntry)
	vbox.Add(endTimeEntry)

	facilityNames = make(map[string]int)
	updateFacilities := func() {
		if validateDate(startDateEntry.Text) && validateTime(startTimeEntry.Text) && validateDate(endDateEntry.Text) && validateTime(endTimeEntry.Text) {
			facilities, err := e.facilityServ.GetFacilitiesByDate(startDateEntry.Text, startTimeEntry.Text, endDateEntry.Text, endTimeEntry.Text, 0, 0)
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

	startDateEntry.OnChanged = func(string) { updateFacilities() }
	startTimeEntry.OnChanged = func(string) { updateFacilities() }
	endDateEntry.OnChanged = func(string) { updateFacilities() }
	endTimeEntry.OnChanged = func(string) { updateFacilities() }
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
