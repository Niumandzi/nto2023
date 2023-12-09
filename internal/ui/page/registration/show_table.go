package registration

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (r RegistrationPage) ShowTableRegistration(window fyne.Window) fyne.CanvasObject {
	registrations, err := r.registrationServ.GetRegistrations(0, 0, 0)
	if err != nil {
		dialog.ShowError(err, window)
		return nil
	}

	dayTranslations := map[string]string{
		"monday":    "Понедельник",
		"tuesday":   "Вторник",
		"wednesday": "Среда",
		"thursday":  "Четверг",
		"friday":    "Пятница",
		"saturday":  "Суббота",
		"sunday":    "Воскресенье",
	}

	daysOfWeek := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	table := container.NewGridWithColumns(len(daysOfWeek) + 1)

	table.Add(widget.NewLabel(""))
	for _, day := range daysOfWeek {
		table.Add(widget.NewLabel(dayTranslations[day]))
	}

	for _, registration := range registrations {
		table.Add(widget.NewLabel(registration.Name))

		for _, day := range daysOfWeek {
			var cellContent string
			for _, schedule := range registration.Schedule {
				if schedule.Day == day {
					cellContent = fmt.Sprintf("%s - %s\n%s\n%s", schedule.StartTime, schedule.EndTime, registration.Facility.Name, registration.Teacher.Name)
					break
				}
			}

			table.Add(widget.NewLabel(cellContent))
		}
	}

	scroll := container.NewScroll(table)

	return scroll
}
