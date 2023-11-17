package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	detailsRepository "github.com/niumandzi/nto2023/internal/repository/details"
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
	detailsRepo := detailsRepository.NewDetailsRepository(db, logger)
	eventServ := eventService.NewEventService(eventRepo, detailsRepo, timeoutContext, logger, ctx)
	event := eventPage.NewEventPage(eventServ, logger)

	//id, err := eventServ.CreateDetails("entertainment", "театр")
	//id, err = eventServ.CreateDetails("entertainment", "кино")
	//
	//teatr1 := model.EventWithDetails{
	//	Id:          0,
	//	Name:        "first teatr",
	//	Description: "pervy",
	//	Date:        "1",
	//	Details:     model.Details{Category: "entertainment", TypeName: "театр"},
	//}
	//id, err = eventServ.CreateEvent(teatr1)
	//
	//teatr2 := model.EventWithDetails{
	//	Id:          0,
	//	Name:        "second teatr",
	//	Description: "vtoroy",
	//	Date:        "2",
	//	Details:     model.Details{Category: "entertainment", TypeName: "театр"},
	//}
	//id, err = eventServ.CreateEvent(teatr2)
	//
	//teatr3 := model.EventWithDetails{
	//	Id:          0,
	//	Name:        "third teatr",
	//	Description: "trety",
	//	Date:        "3",
	//	Details:     model.Details{Category: "entertainment", TypeName: "кино"},
	//}
	//id, err = eventServ.CreateEvent(teatr3)
	//
	//data, _ := eventServ.GetDetailsByCategory("entertainment")
	//fmt.Printf("details by category: %v", data)
	//
	//data2, _ := eventServ.GetEventsByCategory("entertainment")
	//fmt.Printf("events by category: %v", data2)
	//
	//data3, _ := eventServ.GetEventsByCategoryAndType("entertainment", "театр")
	//fmt.Printf("events by category and type: %v", data3)
	//_, _ = id, err

	gui := ui.NewGUI(a, w, event)
	ui.SetupUI(gui)
}
