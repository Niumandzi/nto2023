package registration

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
	"strings"
)

func (r RegistrationPage) ShowRegistration(facilityID int, mugID int, teacherID int, window fyne.Window, bookingContainer *fyne.Container) {
	registrations, err := r.registrationServ.GetRegistrations(facilityID, mugID, teacherID)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	bookingContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, registration := range registrations {
		card := r.createRegistrationCard(registration, window, func() {
			r.ShowRegistration(facilityID, mugID, teacherID, window, bookingContainer)
		})
		grid.Add(card)
	}

	bookingContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	bookingContainer.Refresh()
}

func (r RegistrationPage) createRegistrationCard(registration model.RegistrationWithDetails, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := combineCards(registration)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		//bookingToUpdate := model.BookingWithFacility{
		//	ID:          booking.ID,
		//	Description: booking.Description,
		//	CreateDate:  booking.CreateDate,
		//	StartDate:   booking.StartDate,
		//	StartTime:   booking.StartTime,
		//	EndDate:     booking.EndDate,
		//	EndTime:     booking.EndTime,
		//	Event:       booking.Event,
		//	Facility:    booking.Facility,
		//	Parts:       booking.Parts,
		//}
		//b.UpdateBooking(categoryName, bookingToUpdate, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := r.registrationServ.DeleteRegistration(registration.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Кружок удален", "Кружок успешно удален!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)

	bookingContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return bookingContainer
}

func combineCards(registration model.RegistrationWithDetails) string {
	return registrationCard(registration) + "\n\n" +
		scheduleCard(registration) + "\n\n" +
		facilityCard(registration)
}

func registrationCard(registration model.RegistrationWithDetails) string {
	return fmt.Sprintf("Название: %s\nТип: %s\nПреподаватель: %s\nДата начала: %s",
		registration.Name, registration.MugType.Name, registration.Teacher.Name, registration.StartDate)
}

func scheduleCard(registration model.RegistrationWithDetails) string {

	result := fmt.Sprintf("Количество дней в неделю: %d", registration.NumberOfDays)

	if len(registration.Schedule) > 0 {
		var partsInfo []string
		var day string

		for _, schedule := range registration.Schedule {
			switch schedule.Day {
			case "monday":
				day = "Понедельник"
			case "tuesday":
				day = "Вторник"
			case "wednesday":
				day = "Среда"
			case "thursday":
				day = "Четверг"
			case "friday":
				day = "Пятница"
			case "saturday":
				day = "Суббота"
			case "sunday":
				day = "Воскресенье"
			}

			partsInfo = append(partsInfo, fmt.Sprintf("День недели: %s\nВремя начала: %s\nВремя окончания: %s", day, schedule.StartTime, schedule.EndTime))
		}
		result += "\n" + strings.Join(partsInfo, "\n")
	}

	return result
}

func facilityCard(registration model.RegistrationWithDetails) string {
	result := fmt.Sprintf("Помещение: %s", registration.Facility.Name)

	if len(registration.Parts) > 0 {
		var partsInfo []string
		for _, part := range registration.Parts {
			partsInfo = append(partsInfo, fmt.Sprintf("Часть: %s", part.Name))
		}
		result += "\n" + strings.Join(partsInfo, "\n")
	} else {
		result += "\nЧасти: нет"
	}

	return result
}
