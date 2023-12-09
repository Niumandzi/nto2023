package registration

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/model"
)

func (r RegistrationPage) CreateRegistration(window fyne.Window, onUpdate func()) {
	vbox := container.NewVBox()

	var selectedMugTypeID int
	var selectedTeacherID int
	var selectedFacilityID int
	var selectedNumberOfDays int
	var facilityNames map[string]int
	var facilityParts map[int]map[int]string
	var facilitySelect *widget.Select

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

	numberOfDays := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	daysOfWeek := map[string]bool{"Понедельник": false, "Вторник": false, "Среда": false,
		"Четверг": false, "friday": false, "Суббота": false, "Воскресенье": false}

	dayTimeEntries := make(map[string][2]*widget.Entry)
	daysBox := container.NewVBox()

	dayContainers := make(map[string]*fyne.Container)

	var schedule []model.Schedule

	for day := range daysOfWeek {
		localDay := day

		updateSchedule := func() {
			for i, scheduleRecord := range schedule {
				if scheduleRecord.Day == localDay {
					entryPair, exists := dayTimeEntries[localDay]
					if exists {
						schedule[i].StartTime = entryPair[0].Text
						schedule[i].EndTime = entryPair[1].Text
					}
					break
				}
			}
		}

		dayCheck := widget.NewCheck(localDay, func(checked bool) {
			daysOfWeek[localDay] = checked
			if checked {
				startTimeEntry := widget.NewEntry()
				startTimeEntry.OnChanged = func(string) { updateSchedule() }
				endTimeEntry := widget.NewEntry()
				endTimeEntry.OnChanged = func(string) { updateSchedule() }

				dayTimeEntries[localDay] = [2]*widget.Entry{startTimeEntry, endTimeEntry}
				dayContainer := container.NewHBox(widget.NewLabel(localDay), startTimeEntry, endTimeEntry)
				dayContainers[localDay] = dayContainer
				daysBox.Add(dayContainer)

				schedule = append(schedule, model.Schedule{
					Day:       "friday",
					StartTime: "15:00",
					EndTime:   "18:00",
				})
			} else {
				if dayContainer, exists := dayContainers[localDay]; exists {
					daysBox.Remove(dayContainer)
					delete(dayTimeEntries, localDay)
					delete(dayContainers, localDay)

					for i, scheduleRecord := range schedule {
						if scheduleRecord.Day == localDay {
							schedule = append(schedule[:i], schedule[i+1:]...)
							break
						}
					}
				}
			}
			window.Canvas().Refresh(vbox)
		})

		daysBox.Add(dayCheck)
	}

	numberOfDaysSelect := component.SelectorWidget("Количество занятий в неделю", numberOfDays, func(id int) {
		selectedNumberOfDays = id
		vbox.Add(daysBox)
		window.Canvas().Refresh(vbox)
	}, nil)
	nameEntry := component.EntryWidget("Название")
	startDateEntry := component.EntryWidget("Дата начала (гггг-мм-дд)")

	var selectedParts []int

	var partsBox *fyne.Container
	partsBox = container.NewVBox()
	vbox.Add(partsBox)

	var customPopUp *widget.PopUp

	saveButton := widget.NewButton("            Создать            ", func() {

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
	vbox.Add(numberOfDaysSelect)

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
		vbox.Remove(facilitySelect)
		vbox.Remove(buttons)
		facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
			selectedFacilityID = id
			updateParts()
		}, nil)
		vbox.Add(facilitySelect)
		vbox.Add(buttons)
	}

	facilitySelect = component.SelectorWidget("Помещение", facilityNames, func(id int) {
		selectedFacilityID = id
		updateParts()
	}, nil)

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
		dialog.ShowInformation("Бронирование создано", "Бронирование успешно создано!", window)
		onUpdate()
	}
}
