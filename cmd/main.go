package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	detailsRepository "github.com/niumandzi/nto2023/internal/repository/details"
	eventRepository "github.com/niumandzi/nto2023/internal/repository/event"
	eventService "github.com/niumandzi/nto2023/internal/service/event"
	"github.com/niumandzi/nto2023/internal/ui"
	eventPage "github.com/niumandzi/nto2023/internal/ui/page/event"
	"github.com/niumandzi/nto2023/pkg/logging"
	"github.com/niumandzi/nto2023/pkg/sqlitedb"
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
	eventServ := eventService.NewEventService(eventRepo, detailsRepo, timeoutContext, logger, ctx)
	event := eventPage.NewEventPage(eventServ, logger)

	gui := ui.NewGUI(a, w, event)
	ui.SetupUI(gui)
}
