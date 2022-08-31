package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func GetContent() *fyne.Container {
	mainContent := container.New(layout.NewHBoxLayout())
	return mainContent
}

func GetSettingsPage() *fyne.Container {
	return container.New(layout.NewVBoxLayout())
}
