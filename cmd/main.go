package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	detailsRepository "github.com/niumandzi/nto2023/internal/repository/details"
	eventRepository "github.com/niumandzi/nto2023/internal/repository/event"
	facilityRepository "github.com/niumandzi/nto2023/internal/repository/facility"
	workTypeRepository "github.com/niumandzi/nto2023/internal/repository/work"
	detailsService "github.com/niumandzi/nto2023/internal/service/details"
	eventService "github.com/niumandzi/nto2023/internal/service/event"
	facilityService "github.com/niumandzi/nto2023/internal/service/facility"
	workTypeService "github.com/niumandzi/nto2023/internal/service/work"
	"github.com/niumandzi/nto2023/internal/ui"
	detailsPage "github.com/niumandzi/nto2023/internal/ui/page/details"
	eventPage "github.com/niumandzi/nto2023/internal/ui/page/event"
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

	eventRepo := eventRepository.NewEventRepository(db, logger)
	detailsRepo := detailsRepository.NewDetailsRepository(db, logger)
	facilityRepo := facilityRepository.NewFacilityRepository(db, logger)
	workTypeRepo := workTypeRepository.NewWorkTypeRepository(db, logger)

	eventServ := eventService.NewEventService(eventRepo, detailsRepo, timeoutContext, logger, ctx)
	detailsServ := detailsService.NewDetailsService(detailsRepo, timeoutContext, logger, ctx)
	facilityService.NewFacilityService(facilityRepo, timeoutContext, logger, ctx)
	workTypeServ := workTypeService.NewWorkTypeService(workTypeRepo, timeoutContext, logger, ctx)

	event := eventPage.NewEventPage(eventServ, logger)
	details := detailsPage.NewDetailsPage(detailsServ, logger)

	workType := workPage.NewWorkTypePage(workTypeServ, logger)

	gui := ui.NewGUI(a, w)
	ui.SetupUI(gui, event, details, workType)
}
