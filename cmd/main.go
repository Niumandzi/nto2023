package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	eventRepository "github.com/niumandzi/nto2023/internal/repository/event"
	eventService "github.com/niumandzi/nto2023/internal/service/event"
	"github.com/niumandzi/nto2023/internal/ui"
	"github.com/niumandzi/nto2023/internal/ui/component"
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
	w := a.NewWindow("NTO 2023")
	w.Resize(fyne.NewSize(1200, 700))

	db, err := sqlitedb.NewClient("sqlite3", "./nto2023.db")
	if err != nil {
		component.ShowErrorDialogWidget(err, w)
		logger.Errorf(err.Error())
	}

	err = sqlitedb.CreateTables(db)
	if err != nil {
		component.ShowErrorDialogWidget(err, w)
		logger.Errorf(err.Error())
	}

	timeoutContext := time.Duration(2) * time.Second

	eventRepo := eventRepository.NewEventRepository(db, logger)
	eventServ := eventService.NewEventService(eventRepo, timeoutContext, logger, ctx)
	event := eventPage.NewEventPage(eventServ, logger)

	gui := ui.NewGUI(a, w, event)
	ui.SetupUI(gui)
}
