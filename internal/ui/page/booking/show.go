package booking

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

func (s BookingPage) ShowBooking(startDate string, endDate string, eventID int, categoryName string, window fyne.Window, bookingContainer *fyne.Container) {
	bookings, err := s.bookingServ.GetBookings(startDate, endDate, eventID, categoryName)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	bookingContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, booking := range bookings {
		card := s.createBookingCard(booking, categoryName, window, func() {
			s.ShowBooking(startDate, endDate, eventID, categoryName, window, bookingContainer)
		})
		grid.Add(card)
	}

	bookingContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	bookingContainer.Refresh()
}

func (s BookingPage) createBookingCard(booking model.BookingWithFacility, categoryName string, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := combineCards(booking, categoryName)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		//bookingToUpdate := model.Booking{
		//	ID:          booking.ID,
		//	Description: booking.Description,
		//	CreateDate:  booking.CreateDate,
		//	StartDate:   booking.StartDate,
		//	EndDate:     booking.EndDate,
		//	EventID:     booking.Event.ID,
		//	FacilityID:  booking.Facility.ID,
		//	//PartIDs:     booking.,
		//}
		//s.UpdateBooking(booking.Details.Category, booking.Details.TypeName, bookingToUpdate, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.bookingServ.DeleteBooking(booking.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Событие удалено", "Событие успешно удалено!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)

	bookingContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return bookingContainer
}

func combineCards(booking model.BookingWithFacility, categoryName string) string {
	return eventCard(booking, categoryName) + "\n" + facilityCard(booking) + "\n" + cardBooking(booking)
}

func eventCard(booking model.BookingWithFacility, categoryName string) string {
	var category string
	var categoryLine string

	switch booking.Event.Details.TypeName {
	case "entertainment":
		category = "Развлечения"
	case "enlightenment":
		category = "Просвещение"
	case "education":
		category = "Образование"
	}

	if categoryName == "" {
		categoryLine = fmt.Sprintf("Тип: %s\n", category)
	}

	return categoryLine + fmt.Sprintf("Название: %s\nДата: %s\nОписание: %s",
		booking.Event.Name, booking.Event.Date, booking.Event.Description)
}

func facilityCard(booking model.BookingWithFacility) string {
	result := fmt.Sprintf("Помещение: %s", booking.Facility.Name)

	if len(booking.Parts) > 0 {
		partsInfo := []string{}
		for _, part := range booking.Parts {
			partsInfo = append(partsInfo, fmt.Sprintf("Часть: %s", part.Name))
		}
		result += "\n" + strings.Join(partsInfo, "\n")
	} else {
		result += "\nЧасти: нет"
	}

	return result
}

func cardBooking(booking model.BookingWithFacility) string {
	return fmt.Sprintf("Описание: %s\nДата создания: %s\nДата начала: %s\nДата окончания: %s",
		booking.Description, booking.CreateDate, booking.StartDate, booking.EndDate)
}
