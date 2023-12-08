package registration

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/niumandzi/nto2023/internal/service"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type RegistrationPage struct {
	registrationServ service.RegistrationService
	logger           logging.Logger
}

func NewRegistrationPage(reg service.RegistrationService, logger logging.Logger) RegistrationPage {
	return RegistrationPage{
		registrationServ: reg,
		logger:           logger,
	}
}

func (r RegistrationPage) IndexRegistration(window fyne.Window) fyne.CanvasObject {
	registrationContainer := container.NewStack()
	registrationList := func(facilityID int, mugID int, teacherID int) {
		r.ShowRegistration(facilityID, mugID, teacherID, window, registrationContainer)
	}

	//typeSelect := component.SelectorWidget("Тип мероприятия", typeNames, func(id int) {
	//	bookingList("", "", 0, categoryName)
	//},
	//	nil,
	//)

	createBookingButton := widget.NewButton("Создать бронь", func() {
		//r.CreateRegistration(window, func() {
		//	registrationList(0, 0, 0)
		//})
	})
	createButtons := container.NewHBox(createBookingButton)

	toolbar := container.NewBorder(nil, nil, nil, createButtons)
	content := container.NewBorder(toolbar, nil, nil, nil, registrationContainer)
	registrationList(0, 0, 0)

	return content
}
