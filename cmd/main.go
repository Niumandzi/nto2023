package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	teacherPage "github.com/niumandzi/nto2023/internal/ui/page/teacher"

	applicationRepository "github.com/niumandzi/nto2023/internal/repository/application"
	bookingRepository "github.com/niumandzi/nto2023/internal/repository/booking"
	detailsRepository "github.com/niumandzi/nto2023/internal/repository/details"
	eventRepository "github.com/niumandzi/nto2023/internal/repository/event"
	facilityRepository "github.com/niumandzi/nto2023/internal/repository/facility"
	mugTypeRepository "github.com/niumandzi/nto2023/internal/repository/mug"
	partRepository "github.com/niumandzi/nto2023/internal/repository/part"
	teacherRepository "github.com/niumandzi/nto2023/internal/repository/teacher"
	workTypeRepository "github.com/niumandzi/nto2023/internal/repository/work"

	applicationService "github.com/niumandzi/nto2023/internal/service/application"
	bookingService "github.com/niumandzi/nto2023/internal/service/booking"
	detailsService "github.com/niumandzi/nto2023/internal/service/details"
	eventService "github.com/niumandzi/nto2023/internal/service/event"
	facilityService "github.com/niumandzi/nto2023/internal/service/facility"
	mugTypeService "github.com/niumandzi/nto2023/internal/service/mug"
	partService "github.com/niumandzi/nto2023/internal/service/part"
	teacherService "github.com/niumandzi/nto2023/internal/service/teacher"
	workTypeService "github.com/niumandzi/nto2023/internal/service/work"

	"github.com/niumandzi/nto2023/internal/ui"

	applicationPage "github.com/niumandzi/nto2023/internal/ui/page/application"
	bookingPage "github.com/niumandzi/nto2023/internal/ui/page/booking"
	detailsPage "github.com/niumandzi/nto2023/internal/ui/page/details"
	eventPage "github.com/niumandzi/nto2023/internal/ui/page/event"
	facilityPage "github.com/niumandzi/nto2023/internal/ui/page/facility"
	mugPage "github.com/niumandzi/nto2023/internal/ui/page/mug"
	workPage "github.com/niumandzi/nto2023/internal/ui/page/work"

	"github.com/niumandzi/nto2023/pkg/logging"
	"github.com/niumandzi/nto2023/pkg/sqlitedb"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	a := app.New()
	w := a.NewWindow("НТО 2023")
	w.Resize(fyne.NewSize(1200, 700))

	iconBytes, err := os.ReadFile("icon.png")
	if err != nil {
		dialog.ShowError(err, w)
	}

	iconResource := fyne.NewStaticResource("IconName", iconBytes)
	w.SetIcon(iconResource)

	db, err := sqlitedb.NewClient("sqlite3", "./nto2023.db")
	if err != nil {
		dialog.ShowError(err, w)
		logger.Error(err.Error())
	}

	err = sqlitedb.CreateTables(db)
	if err != nil {
		dialog.ShowError(err, w)
		logger.Error(err.Error())
	}

	timeoutContext := time.Duration(2) * time.Second

	applicationRepo := applicationRepository.NewApplicationRepository(db, logger)
	bookingRepo := bookingRepository.NewBookingRepository(db, logger)
	detailsRepo := detailsRepository.NewDetailsRepository(db, logger)
	eventRepo := eventRepository.NewEventRepository(db, logger)
	facilityRepo := facilityRepository.NewFacilityRepository(db, logger)
	mugTypeRepo := mugTypeRepository.NemMugTypeRepository(db, logger)
	partRepo := partRepository.NewPartRepository(db, logger)
	teacherRepo := teacherRepository.NewTeacherRepository(db, logger)
	workTypeRepo := workTypeRepository.NewWorkTypeRepository(db, logger)

	applicationServ := applicationService.NewApplicationService(applicationRepo, timeoutContext, logger, ctx)
	bookingServ := bookingService.NewBookingService(bookingRepo, timeoutContext, logger, ctx)
	detailsServ := detailsService.NewDetailsService(detailsRepo, timeoutContext, logger, ctx)
	eventServ := eventService.NewEventService(eventRepo, detailsRepo, timeoutContext, logger, ctx)
	facilityServ := facilityService.NewFacilityService(facilityRepo, timeoutContext, logger, ctx)
	mugTypeServ := mugTypeService.NewMugTypeService(mugTypeRepo, timeoutContext, logger, ctx)
	partServ := partService.NewPartService(partRepo, timeoutContext, logger, ctx)
	teacherServ := teacherService.NewTeacherService(teacherRepo, timeoutContext, logger, ctx)
	workTypeServ := workTypeService.NewWorkTypeService(workTypeRepo, timeoutContext, logger, ctx)

	application := applicationPage.NewApplicationPage(applicationServ, eventServ, facilityServ, workTypeServ, logger)
	booking := bookingPage.NewBookingPage(bookingServ, eventServ, facilityServ, logger)
	details := detailsPage.NewDetailsPage(detailsServ, logger)
	event := eventPage.NewEventPage(facilityServ, bookingServ, eventServ, detailsServ, logger)
	facility := facilityPage.NewFacilityPage(facilityServ, partServ, logger)
	mugType := mugPage.NewMugTypePage(mugTypeServ, logger)
	teacher := teacherPage.NewTeacherPage(teacherServ, logger)
	workType := workPage.NewWorkTypePage(workTypeServ, logger)

	gui := ui.NewGUI(a, w)
	ui.SetupUI(gui, application, booking, details, event, facility, mugType, teacher, workType)
}
