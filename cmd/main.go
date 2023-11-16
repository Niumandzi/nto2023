package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/niumandzi/nto2023/internal/ui"
	"github.com/niumandzi/nto2023/internal/ui/component"
	"github.com/niumandzi/nto2023/pkg/logging"
	"github.com/niumandzi/nto2023/pkg/sqlitedb"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	a := app.New()
	w := a.NewWindow("NTO 2023")
	w.Resize(fyne.NewSize(1200, 700))

	db, err := sqlitedb.NewClient("sqlite3", "./nto2022.db")
	if err != nil {
		component.ShowErrorDialog(err, w)
		logger.Errorf(err.Error())
	}

	err = sqlitedb.CreateTables(db)
	if err != nil {
		component.ShowErrorDialog(err, w)
		logger.Errorf(err.Error())
	}

	gui := ui.NewGUI(a, w)
	ui.SetupUI(gui)
}
