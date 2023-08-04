package ui

import "fyne.io/fyne/v2/container"

func Setup(app *AppInit) {
	SetupMenus(app)
	swatchesContainer := BuildSwatches(app)
	colorPicker := SetUpColorPicker(app)
	//NewBorder(top border row, bottom border row, left column border, right border column, center canvas)
	appLayout := container.NewBorder(nil, swatchesContainer, nil, colorPicker, app.PixlCanvas)
	app.PixlWindow.SetContent(appLayout)
}
