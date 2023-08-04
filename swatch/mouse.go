package swatch

import (
	"fyne.io/fyne/v2/driver/desktop"
)

func (swatch *Swatch) MouseDown(ev *desktop.MouseEvent) {
	//select swatch
	swatch.clickHandler(swatch)
	swatch.Selected = true
	swatch.Refresh()
}

func (swatch *Swatch) MouseUp(ev *desktop.MouseEvent) {}
