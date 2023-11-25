package application

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationPage) ShowApplication(categoryName string, facilityId int, workTypeId int, status string, window fyne.Window, applicationContainer *fyne.Container) {
	applications, err := s.applicationServ.GetApplications(categoryName, facilityId, workTypeId, status)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	applicationContainer.Objects = nil

	grid := container.New(layout.NewGridLayoutWithColumns(2))
	for _, application := range applications {
		card := s.createApplicationCard(application, window, func() {
			s.ShowApplication(categoryName, facilityId, workTypeId, status, window, applicationContainer)
		})
		grid.Add(card)
	}

	applicationContainer.Objects = []fyne.CanvasObject{container.NewVScroll(grid)}
	applicationContainer.Refresh()
}

func (s ApplicationPage) createApplicationCard(application model.ApplicationWithDetails, window fyne.Window, onUpdate func()) fyne.CanvasObject {
	cardText := card(application)
	label := widget.NewLabel(cardText)
	label.Wrapping = fyne.TextWrapWord

	updateButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		//s.UpdateApplication(application., application.ID, application.Name, application.Date, application.Description, application.Details.ID, window, onUpdate)
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
	applicationContainer := widget.NewCard("", "", container.NewBorder(nil, buttons, nil, nil, label))

	return applicationContainer
}

func card(application model.ApplicationWithDetails) string {
	return fmt.Sprintf("Тип работ: %s\nПомещение: %s\nОписание: %s\nДата создания: %s\nДата выполнения: %s\nСтатус: %s\nТип: %s\nНазвание: %s\nДата: %s\nОписание: %s",
		application.WorkType.Name, application.Facility.Name, application.Description, application.CreateDate, application.DueDate, application.Status, application.Event.Details.TypeName, application.Event.Name, application.Event.Date, application.Event.Description)
}