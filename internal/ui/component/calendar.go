package component

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func CalendarWidget(parent fyne.Window) *widget.PopUp {
	yearSelect := widget.NewSelect(getYears(), nil)
	monthSelect := widget.NewSelect(getMonths(), nil)

	selectedYear := time.Now().Year()
	selectedMonth := int(time.Now().Month())
	selectedDay := 1

	var popUp *widget.PopUp

	updateDays := func() {
		year, _ := strconv.Atoi(yearSelect.Selected)
		month := time.Month(monthSelect.SelectedIndex() + 1)
		daysInMonth := getDaysInMonth(year, month)

		dayButtons := make([]fyne.CanvasObject, daysInMonth)
		for i := 1; i <= daysInMonth; i++ {
			day := i
			dayButtons[i-1] = widget.NewButton(fmt.Sprintf("%d", i), func() {
				selectedDay = day
				fmt.Printf("Выбранная дата: %d-%02d-%02d\n", selectedYear, selectedMonth, selectedDay)
				popUp.Hide()
			})
		}

		dayContainer := container.NewGridWithColumns(7, dayButtons...)
		popUp.Content = container.NewVBox(yearSelect, monthSelect, dayContainer)
		popUp.Refresh()
	}

	monthSelect.OnChanged = func(string) {
		selectedMonth = monthSelect.SelectedIndex() + 1
		updateDays()
	}

	updateDays()

	popUp = widget.NewModalPopUp(container.NewVBox(yearSelect, monthSelect), parent.Canvas())
	return popUp
}

func getYears() []string {
	currentYear := time.Now().Year()
	years := make([]string, 10)
	for i := 0; i < 10; i++ {
		years[i] = fmt.Sprintf("%d", currentYear+i)
	}
	return years
}

func getMonths() []string {
	months := make([]string, 12)
	for i := 1; i <= 12; i++ {
		months[i-1] = time.Month(i).String()
	}
	return months
}

func getDaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
