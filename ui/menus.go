package ui

import (
	"errors"
	"image"
	"image/png"
	"os"
	"pixl/util"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func SaveFileDialog(app *AppInit) {
	//save file popup
	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, e error) {
		if uri == nil {
			//if no file path, close
			return
		} else {
			err := png.Encode(uri, app.PixlCanvas.PixelData)
			if err != nil {
				dialog.ShowError(err, app.PixlWindow)
				return
			}
			//after user selects file path, file will save to the correct location in future
			app.State.SetFilePath(uri.URI().Path())
		}

	}, app.PixlWindow)
}

func BuildSaveAsMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save As...", func() {
		SaveFileDialog(app)
	})
}

func BuildSaveMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save", func() {
		if app.State.FilePath == "" {
			//if no file save path specified, launch save as func
			SaveFileDialog(app)
		} else {
			tryClose := func(fh *os.File) {
				err := fh.Close()
				if err != nil {
					dialog.ShowError(err, app.PixlWindow)
				}
			}
			fh, err := os.Create(app.State.FilePath)
			defer tryClose(fh)

			if err != nil {
				dialog.ShowError(err, app.PixlWindow)
				return
			}
			err = png.Encode(fh, app.PixlCanvas.PixelData)
			if err != nil {
				dialog.ShowError(err, app.PixlWindow)
				return
			}
		}

	})
}

func BuildNewMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("New", func() {
		sizeValidator := func(s string) error {
			width, err := strconv.Atoi(s)
			if err != nil {
				return errors.New("Must be a positive integer")
			}
			if width <= 0 {
				return errors.New("Must have width > 0")
			}
			return nil
		}
		//text entry fields
		widthEntry := widget.NewEntry()
		widthEntry.Validator = sizeValidator

		heightEntry := widget.NewEntry()
		heightEntry.Validator = sizeValidator

		widthFormEntry := widget.NewFormItem("Width", widthEntry)
		heightFormEntry := widget.NewFormItem("Height", heightEntry)

		formItems := []*widget.FormItem{widthFormEntry, heightFormEntry}

		dialog.ShowForm("New Iage", "Create", "Cancel", formItems, func(ok bool) {
			if ok {
				pixelWidth := 0
				pixelHeight := 0
				if widthEntry.Validate() != nil {
					dialog.ShowError(errors.New("Invalid Width"), app.PixlWindow)
				} else {
					pixelWidth, _ = strconv.Atoi(widthEntry.Text)
				}

				if heightEntry.Validate() != nil {
					dialog.ShowError(errors.New("Invalid Height"), app.PixlWindow)
				} else {
					pixelHeight, _ = strconv.Atoi(heightEntry.Text)
				}

				app.PixlCanvas.NewDrawing(pixelWidth, pixelHeight)

			}
		}, app.PixlWindow)
	})

}

func BuildOpenMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Open...", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, e error) {
			if uri == nil {
				return
			} else {
				//open image
				image, _, err := image.Decode(uri)
				if err != nil {
					dialog.ShowError(err, app.PixlWindow)
					return
				}
				app.PixlCanvas.LoadImage(image)
				app.State.SetFilePath(uri.URI().Path())
				//run image through color picker to log each color in each pixel
				imgColors := util.GetImageColors(image)
				//then populate each canvas pixel with those colors
				i := 0
				for c := range imgColors {
					if i == len(app.Swatches) {
						break
					}
					app.Swatches[i].SetColor(c)
					i++
				}
			}
		}, app.PixlWindow)
	})
}

func BuildMenus(app *AppInit) *fyne.Menu {
	return fyne.NewMenu(
		"File",
		BuildNewMenu(app),
		BuildOpenMenu(app),
		BuildSaveMenu(app),
		BuildSaveAsMenu(app),
	)
}

func SetupMenus(app *AppInit) {
	menus := BuildMenus(app)
	mainMenu := fyne.NewMainMenu(menus)
	app.PixlWindow.SetMainMenu(mainMenu)
}
