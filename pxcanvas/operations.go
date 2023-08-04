package pxcanvas

import "fyne.io/fyne/v2"

//fyne.PointEvent stores posiitioon of mouse
func (pxCanvas *PxCanvas) Pan(previousCoordinate, currentCoordinate fyne.PointEvent) {
	xDiff := currentCoordinate.Position.X - previousCoordinate.Position.X
	yDiff := currentCoordinate.Position.Y - previousCoordinate.Position.Y
	pxCanvas.CanvasOffset.X += xDiff
	pxCanvas.CanvasOffset.Y += yDiff
	pxCanvas.Refresh()
}

//zoom in or out
func (pxCanvas *PxCanvas) scale(direction int) {
	switch {
	case direction > 0:
		pxCanvas.PxSize++
	case direction < 0:
		if pxCanvas.PxSize > 2 {
			pxCanvas.PxSize--
		}
	default:
		pxCanvas.PxSize = 10
	}
}
