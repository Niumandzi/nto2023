package application

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
	"image/color"
)

func (s ApplicationPage) ShowApplication(categoryName string, facilityId int, workTypeId int, status string, window fyne.Window, applicationContainer *fyne.Container) {
	applications, err := s.applicationServ.GetApplications(categoryName, facilityId, workTypeId, status)
	if err != nil {
		dialog.ShowError(err, window)
	}

	applicationContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(3))
	for _, application := range applications {
		card := s.createApplicationCard(application, categoryName, window, func() {
			s.ShowApplication(categoryName, facilityId, workTypeId, status, window, applicationContainer)
		})
		grid.Add(card)
	}

	applicationContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	applicationContainer.Refresh()
}

func (s ApplicationPage) createApplicationCard(application model.ApplicationWithDetails, categoryName string, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := combineCards(application, categoryName)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		appToUpdate := model.Application{
			ID:          application.ID,
			Description: application.Description,
			CreateDate:  application.CreateDate,
			DueDate:     application.DueDate,
			Status:      application.Status,
			WorkTypeId:  application.WorkType.ID,
			FacilityId:  application.Facility.ID,
			EventId:     application.Event.ID,
		}
		s.UpdateApplication(categoryName, application.WorkType.Name, application.Facility.Name, application.Event.Name, appToUpdate, window, onUpdate)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		err := s.applicationServ.DeleteApplication(application.ID)
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			dialog.ShowInformation("Заявка удалена", "Заявка успешно удалена!", window)
			onUpdate()
		}
	})

	deleteButton.Importance = widget.LowImportance
	updateButton.Importance = widget.LowImportance

	buttons := container.NewHBox(layout.NewSpacer(), updateButton, deleteButton)

	statusColorBar := canvas.NewRectangle(getColorBasedOnStatus(application.Status))
	statusColorBar.SetMinSize(fyne.NewSize(200, 6))

	applicationContainer := widget.NewCard("", "", container.NewBorder(statusColorBar, buttons, nil, nil, label))

	return applicationContainer
}

func combineCards(application model.ApplicationWithDetails, categoryName string) string {
	return applicationCard(application) + "\n" + eventCard(application, categoryName)
}

func applicationCard(application model.ApplicationWithDetails) string {
	var status string

	switch application.Status {
	case "created":
		status = "Черновик"
	case "todo":
		status = "К выполнению"
	case "done":
		status = "Выполнено"
	}

	return fmt.Sprintf("Тип работ: %s\nПомещение: %s\nОписание: %s\nДата создания: %s\nДата выполнения: %s\nСтатус: %s\n",
		application.WorkType.Name, application.Facility.Name, application.Description, application.CreateDate, application.DueDate, status)
}

func eventCard(application model.ApplicationWithDetails, categoryName string) string {
	var category string
	var categoryLine string

	switch application.Event.Details.TypeName {
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
		application.Event.Name, application.Event.Date, application.Event.Description)
}

func getColorBasedOnStatus(status string) color.Color {
	switch status {
	case "todo":
		return color.RGBA{R: 167, G: 130, B: 149, A: 255}
	case "done":
		return color.RGBA{R: 128, G: 128, B: 128, A: 255}
	default:
		return color.Alpha{}
	}
}
