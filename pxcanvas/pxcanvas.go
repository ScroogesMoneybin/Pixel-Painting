package pxcanvas

import (
	"image"
	"image/color"
	"pixl/apptype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	DesignatedScrollButton = desktop.RightMouseButton
)

type PxCanvasMouseState struct {
	previousCoord *fyne.PointEvent
}

type PxCanvas struct {
	widget.BaseWidget
	apptype.PxCanvasConfig
	renderer    *PxCanvasRenderer
	PixelData   image.Image
	mouseState  PxCanvasMouseState
	appState    *apptype.State
	reloadImage bool
}

func (pxCanvas *PxCanvas) Bounds() image.Rectangle {
	//x0 and y0 is top left corner of canvas
	x0 := int(pxCanvas.CanvasOffset.X)
	y0 := int(pxCanvas.CanvasOffset.Y)
	//x1 and y1 is coordinate of bottom right corner of canvas
	x1 := int(pxCanvas.PxCols*pxCanvas.PxSize + int(pxCanvas.CanvasOffset.X))
	y1 := int(pxCanvas.PxRows*pxCanvas.PxSize + int(pxCanvas.CanvasOffset.Y))
	return image.Rect(x0, y0, x1, y1)
}

// determine if mouse is on canvas
func InBounds(pos fyne.Position, bounds image.Rectangle) bool {
	if pos.X >= float32(bounds.Min.X) && pos.X < float32(bounds.Max.X) && pos.Y >= float32(bounds.Min.Y) && pos.Y < float32(bounds.Max.Y) {
		return true
	}
	return false
}

func NewBlankImage(cols, rows int, c color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, cols, rows))
	//go thru each image pixel and set the color
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			img.Set(x, y, c)
		}

	}
	return img
}

// make new pixel canvas
func NewPixelCanvas(state *apptype.State, config apptype.PxCanvasConfig) *PxCanvas {
	pxCanvas := &PxCanvas{
		PxCanvasConfig: config,
		appState:       state,
	}
	pxCanvas.PixelData = NewBlankImage(config.PxCols, config.PxRows, color.NRGBA{128, 128, 128, 128})
	pxCanvas.ExtendBaseWidget(pxCanvas) /*provides defaults for fyne toolkit like sizing and position information*/
	return pxCanvas
}

func (pxCanvas *PxCanvas) CreateRenderer() fyne.WidgetRenderer {
	canvasImage := canvas.NewImageFromImage(pxCanvas.PixelData)
	canvasImage.ScaleMode = canvas.ImageScalePixels
	canvasImage.FillMode = canvas.ImageFillContain

	canvasBorder := make([]canvas.Line, 4)
	//defaults for borders
	for i := 0; i < len(canvasBorder); i++ {
		canvasBorder[i].StrokeColor = color.NRGBA{100, 100, 100, 255}
		canvasBorder[i].StrokeWidth = 2
	}
	renderer := &PxCanvasRenderer{
		pxCanvas:     pxCanvas,
		canvasImage:  canvasImage,
		canvasBorder: canvasBorder,
	}
	pxCanvas.renderer = renderer
	return renderer
}

// pan canvas around drawing area
// Mouse event is the pushing of buttons
func (pxCanvas *PxCanvas) TryPan(previousCoord *fyne.PointEvent, ev *desktop.MouseEvent) {
	//pan with right mouse button
	if previousCoord != nil && ev.Button == DesignatedScrollButton {
		pxCanvas.Pan(*previousCoord, ev.PointEvent)
	}
}

// Brushable interface
// access pixel data and set color based on x,y coordinate
func (pxCanvas *PxCanvas) SetColor(c color.Color, x, y int) {
	//check if underlying image is NRGBA or RGBA type, and if so, Set (x,y) coord to color c
	if nrgba, ok := pxCanvas.PixelData.(*image.NRGBA); ok {
		nrgba.Set(x, y, c)
	}
	if rgba, ok := pxCanvas.PixelData.(*image.RGBA); ok {
		rgba.Set(x, y, c)

	}
	pxCanvas.Refresh()
}

func (pxCanvas *PxCanvas) MouseToCanvasXY(ev *desktop.MouseEvent) (*int, *int) {
	bounds := pxCanvas.Bounds()
	if !InBounds(ev.Position, bounds) {
		return nil, nil
	}
	pxSize := float32(pxCanvas.PxSize)
	xOffset := pxCanvas.CanvasOffset.X
	yOffset := pxCanvas.CanvasOffset.Y
	x := int((ev.Position.X - xOffset) / pxSize)
	y := int((ev.Position.Y - yOffset) / pxSize)

	return &x, &y
}

// create new drawing
// load blank and existing images
func (pxCanvas *PxCanvas) LoadImage(img image.Image) {
	dimensions := img.Bounds()
	pxCanvas.PxCanvasConfig.PxCols = dimensions.Dx()
	pxCanvas.PxCanvasConfig.PxRows = dimensions.Dy()

	pxCanvas.PixelData = img
	pxCanvas.reloadImage = true
	pxCanvas.Refresh()
}

func (pxCanvas *PxCanvas) NewDrawing(cols, rows int) {
	pxCanvas.appState.SetFilePath("")
	pxCanvas.PxCols = cols
	pxCanvas.PxRows = rows
	pixelData := NewBlankImage(cols, rows, color.NRGBA{128, 128, 128, 255})
	pxCanvas.LoadImage(pixelData)
}
