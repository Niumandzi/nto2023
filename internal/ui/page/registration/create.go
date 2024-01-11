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

type DayEntry struct {
	DaySelector    *widget.Select
	StartTimeEntry *widget.Entry
	EndTimeEntry   *widget.Entry
	TimeContainer  *fyne.Container
}

func (r RegistrationPage) CreateRegistration(window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()
	daysVBox := container.NewVBox()

	var selectedMugTypeID int
	var selectedTeacherID int
	var selectedFacilityID int
	var selectedNumberOfDays int
	var facilityNames map[string]int
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select
	var dayEntries []DayEntry
	var schedule []model.Schedule
	var customPopUp *widget.PopUp
	var selectedParts []int
	var partsBox *fyne.Container

	mugs, err := r.mugTypeServ.GetActiveMugTypes(0, 0)
	if err != nil {
		dialog.ShowError(err, window)
	}

	mugTypes := make(map[string]int)
	for _, mug := range mugs {
		mugTypes[mug.Name] = mug.ID
	}

	mugTypeSelect := component.SelectorWidget("Тип кружка", mugTypes, func(id int) {
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

	teacherSelect := component.SelectorWidget("Преподаватель", teacherNames, func(id int) {
		selectedTeacherID = id
	}, nil)

	days := map[string]string{
		"Понедельник": "monday",
		"Вторник":     "tuesday",
		"Среда":       "wednesday",
		"Четверг":     "thursday",
		"Пятница":     "friday",
		"Суббота":     "saturday",
		"Воскресенье": "sunday",
	}

	nameEntry := component.EntryWidget("Название")
	startDateEntry := component.EntryWidget("Дата начала (гггг-мм-дд)")

	saveButton := widget.NewButton("            Создать            ", func() {
		if len(schedule) != selectedNumberOfDays {
			dialog.ShowInformation("Предупреждение", "Количество выбранных дней не соответствует выбранному количеству занятий в неделю.", window)
			return
		}

		formData := model.Registration{
			Name:         nameEntry.Text,
			StartDate:    startDateEntry.Text,
			NumberOfDays: selectedNumberOfDays,
			FacilityID:   selectedFacilityID,
			MugTypeID:    selectedMugTypeID,
			TeacherID:    selectedTeacherID,
			Schedule:     schedule,
			PartIDs:      selectedParts,
		}
		handleCreateRegistration(formData, window, r.registrationServ, onUpdate, customPopUp)

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

	addDayButton := widget.NewButton("              Добавить день недели              ", func() {
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
	facilities, err := r.facilityServ.GetActiveFacilities("", 0, "")
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
		vbox.Remove(buttons)
		facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
			selectedFacilityID = id
			updateParts()
		}, nil)
		vbox.Add(buttons)
	}

	facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityID = id
		updateParts()
	}, nil)

	vbox.Add(daysVBox)
	vbox.Add(buttonBox)
	vbox.Add(facilitySelect)
	vbox.Add(buttons)

	customPopUp = widget.NewModalPopUp(vbox, window.Canvas())
	customPopUp.Resize(fyne.NewSize(300, 200))
	customPopUp.Show()
}

func handleCreateRegistration(appData model.Registration, window fyne.Window, registrationServ service.RegistrationService, onUpdate func(), popUp *widget.PopUp) {
	_, err := registrationServ.CreateRegistration(appData)

	if err != nil {
		dialog.ShowError(err, window)
	} else {
		popUp.Hide()
		dialog.ShowInformation("Кружок зарегестрирован", "Кружок успешно зарегестрирован!", window)
		onUpdate()
	}
}
