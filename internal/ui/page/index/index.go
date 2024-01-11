package index

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func ShowIndex() fyne.CanvasObject {
	greetingLabel := widget.NewLabel("Приветствуем!")
	greetingLabel.Alignment = fyne.TextAlignCenter

	descriptionLabel := widget.NewLabel("Это реализация третьего задания НТО 2023, командой 'ЯМы Райн Гослинг'")
	descriptionLabel.Alignment = fyne.TextAlignCenter

	bugLabel := widget.NewLabel("Наше приложение очень удобно логирует все ошибки." + "\n" + "Если вы столкнетесь с какими-либо багами, пожалуйста уведомите нас об этом в тг." + "\n" +
		"Пришлите log файл (лежит в папке рядом с exe), а также видео или скриншот с информацией о том, как вы столкнулись с ошибкой." + "\n" + "Спасибо!")
	bugLabel.Alignment = fyne.TextAlignCenter

	pares, _ := url.Parse("https://t.me/niumandzi")
	hyperlink := widget.NewHyperlink("@niumandzi", pares)
	hyperlink.Alignment = fyne.TextAlignCenter

	// Вертикальное расположение элементов
	content := container.NewVBox(
		greetingLabel,
		descriptionLabel,
		bugLabel,
		hyperlink, // Добавление гиперссылки
	)

	// Центрирование контента
	return container.New(layout.NewCenterLayout(), content)
}
