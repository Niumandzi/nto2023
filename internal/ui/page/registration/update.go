package registration

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (r RegistrationPage) UpdateRegistration(registration model.RegistrationWithDetails, window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()

	selectedMugTypeID := registration.MugType.ID
	selectedTeacherID := registration.Teacher.ID
	selectedFacilityID := registration.Facility.ID
	selectedNumberOfDays := registration.NumberOfDays
	var selectedParts []int
	var facilityNames map[string]int
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select
	var dayEntries []DayEntry
	var schedule []model.Schedule
	var partsBox *fyne.Container

	daysVBox := container.NewVBox()

	facilities, err := r.facilityServ.GetActiveFacilities("", 0, "")
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
	for _, part := range registration.Parts {
		bookingPartIDs[part.ID] = true
	}

	selectedParts = []int{}
	for partID := range bookingPartIDs {
		selectedParts = append(selectedParts, partID)
	}

	newPartsBox := container.NewVBox()
	for partID, partName := range facilityParts[registration.Facility.ID] {
		localPartID := partID

		isChecked := bookingPartIDs[localPartID]

		checkBox := widget.NewCheck(partName, nil)
		checkBox.Checked = isChecked
		checkBox.Refresh()

		checkBox.OnChanged = func(checked bool) {
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
		}
		newPartsBox.Add(checkBox)
	}

	mugs, err := r.mugTypeServ.GetActiveMugTypes(0, 0)
	if err != nil {
		dialog.ShowError(err, window)
	}

	mugTypes := make(map[string]int)
	for _, mug := range mugs {
		mugTypes[mug.Name] = mug.ID
	}

	mugTypeSelect := component.SelectorWidget(registration.MugType.Name, mugTypes, func(id int) {
		selectedMugTypeID = id
	}, nil)

	teachers, err := r.teacherServ.GetActiveTeachers(0, 0)
	if err != nil {
		dialog.ShowError(err, window)
	}

	teacherNames := make(map[string]int)
	for _, teacher := range teachers {
		teacherNames[teacher.Name] = teacher.ID
	}

	teacherSelect := component.SelectorWidget(registration.Teacher.Name, teacherNames, func(id int) {
		selectedTeacherID = id
	}, nil)

	nameEntry := component.EntryWidgetWithData("Название", registration.Name)
	startDateEntry := component.EntryWidgetWithData("Дата начала (гггг-мм-дд)", registration.StartDate)

	partsBox = container.NewVBox()
	vbox.Add(partsBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Сохранить            ", func() {
		if len(schedule) != selectedNumberOfDays {
			dialog.ShowInformation("Предупреждение", "Количество выбранных дней не соответствует выбранному количеству занятий в неделю.", window)
			return
		}

		formData := model.Registration{
			ID:           registration.ID,
			Name:         nameEntry.Text,
			StartDate:    startDateEntry.Text,
			NumberOfDays: selectedNumberOfDays,
			FacilityID:   selectedFacilityID,
			MugTypeID:    selectedMugTypeID,
			TeacherID:    selectedTeacherID,
			Schedule:     schedule,
			PartIDs:      selectedParts,
		}
		handleUpdateRegistration(formData, window, r.registrationServ, onUpdate, customPopUp)

	})
	cancelButton := widget.NewButton("            Отмена            ", func() {
		customPopUp.Hide()
	})

	buttons := container.NewHBox(saveButton, cancelButton)

	updateDaysVBox := func() {
		daysVBox.RemoveAll()
		for _, entry := range dayEntries {
			daysVBox.Add(entry.DaySelector)
			daysVBox.Add(entry.TimeContainer)
		}
		window.Canvas().Refresh(daysVBox)
	}

	days := map[string]string{
		"Понедельник": "monday",
		"Вторник":     "tuesday",
		"Среда":       "wednesday",
		"Четверг":     "thursday",
		"Пятница":     "friday",
		"Суббота":     "saturday",
		"Воскресенье": "sunday",
	}

	daysForTranslate := map[string]string{
		"monday":    "Понедельник",
		"tuesday":   "Вторник",
		"wednesday": "Среда",
		"thursday":  "Четверг",
		"friday":    "Пятница",
		"saturday":  "Суббота",
		"sunday":    "Воскресенье",
	}

	addDayButton := widget.NewButton("                 Добавить день недели                 ", func() {
		if len(dayEntries) >= 3 {
			return
		}

		newDaySelector := component.SelectorWidget("День недели", days, nil, nil)
		newStartTimeEntry := component.EntryWidget("Дата начала")
		newEndTimeEntry := component.EntryWidget("Дата окончания")

		newStartTimeEntry.Disable()
		newEndTimeEntry.Disable()

		newDaySelector.OnChanged = func(selected string) {
			if selected != "" {
				selectedDay := days[selected]

				newDaySelector.Disable()
				newStartTimeEntry.Enable()
				newEndTimeEntry.Enable()

				schedule = append(schedule, model.Schedule{
					Day:       selectedDay,
					StartTime: newStartTimeEntry.Text,
					EndTime:   newEndTimeEntry.Text,
				})
			}
		}

		newStartTimeEntry.OnChanged = func(newTime string) {
			if len(schedule) > 0 && isTimeFormatValid(newTime) {
				schedule[len(schedule)-1].StartTime = newTime
				newStartTimeEntry.Disable()
			}
		}

		newEndTimeEntry.OnChanged = func(newTime string) {
			if len(schedule) > 0 && isTimeFormatValid(newTime) {
				schedule[len(schedule)-1].EndTime = newTime
				newEndTimeEntry.Disable()
			}
		}

		newTimeContainer := container.New(layout.NewGridLayout(2), newStartTimeEntry, newEndTimeEntry)

		newEntry := DayEntry{
			DaySelector:    newDaySelector,
			StartTimeEntry: newStartTimeEntry,
			EndTimeEntry:   newEndTimeEntry,
			TimeContainer:  newTimeContainer,
		}

		selectedNumberOfDays += 1
		dayEntries = append(dayEntries, newEntry)
		updateDaysVBox()
	})

	deleteLastDayButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if len(dayEntries) > 0 {
			lastIndex := len(dayEntries) - 1
			dayEntries = dayEntries[:lastIndex]
			if len(schedule) > lastIndex {
				schedule = schedule[:lastIndex]
			}
			selectedNumberOfDays -= 1
			updateDaysVBox()
		}
	})

	buttonBox := container.NewHBox(addDayButton, deleteLastDayButton)

	createDayEntry := func(day string, startTime string, endTime string) DayEntry {

		newDaySelector := component.SelectorWidget(day, days, nil, nil)
		newStartTimeEntry := component.EntryWidgetWithData("Дата начала", startTime)
		newEndTimeEntry := component.EntryWidgetWithData("Дата окончания", endTime)

		newDaySelector.Disable()
		newStartTimeEntry.Disable()
		newEndTimeEntry.Disable()

		newTimeContainer := container.New(layout.NewGridLayout(2), newStartTimeEntry, newEndTimeEntry)

		return DayEntry{
			DaySelector:    newDaySelector,
			StartTimeEntry: newStartTimeEntry,
			EndTimeEntry:   newEndTimeEntry,
			TimeContainer:  newTimeContainer,
		}
	}

	for _, scheduleItem := range registration.Schedule {
		day := daysForTranslate[scheduleItem.Day]
		newEntry := createDayEntry(day, scheduleItem.StartTime, scheduleItem.EndTime)
		dayEntries = append(dayEntries, newEntry)
		schedule = append(schedule, model.Schedule{
			Day:       scheduleItem.Day,
			StartTime: scheduleItem.StartTime,
			EndTime:   scheduleItem.EndTime,
		})
	}
	updateDaysVBox()

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

	vbox.Add(nameEntry)
	vbox.Add(mugTypeSelect)
	vbox.Add(teacherSelect)
	vbox.Add(startDateEntry)

	facilityNames = make(map[string]int)

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
		vbox.Remove(buttons)
		facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
			selectedFacilityID = id
			updateParts()
		}, nil)
		vbox.Add(buttons)
	}

	facilitySelect = component.SelectorWidget(registration.Facility.Name, facilityNames, func(id int) {
		selectedFacilityID = id
		updateParts()
	}, nil)

	vbox.Add(daysVBox)
	vbox.Add(buttonBox)
	vbox.Add(facilitySelect)
	vbox.Add(buttons)
	partsBox = newPartsBox
	vbox.Add(partsBox)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 200))
	customPopUp.Show()
}

func handleUpdateRegistration(appData model.Registration, window fyne.Window, registrationServ service.RegistrationService, onUpdate func(), popUp *widget.PopUp) {
	err := registrationServ.UpdateRegistration(appData)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		popUp.Hide()
		dialog.ShowInformation("Кружок обновлен", "Кружок успешно обновлен!", window)
		onUpdate()
	}
}
