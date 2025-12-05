package ui

import (
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gridSize            = WindowHeight - 20
	lineWidth           = 12.0
	gridBackgroundColor = "#F4F6F7"
	gridBorderColor     = "#2C3E50"

	cellSize = (gridSize / 3)

	symbolLength = cellSize/2 - 2*lineWidth
)

var (
	GridImage        *ebiten.Image
	CircleImage      *ebiten.Image
	CrossImage       *ebiten.Image
	WheelImage       *ebiten.Image
	MainMenuImage    *ebiten.Image
	WaitingMenuImage *ebiten.Image
	GameMenuImage    *ebiten.Image
)

func InitImages() {
	DrawGrid(GridCol)
	DrawCircle()
	DrawCross()
	DrawWaitingWheel()
	DrawMainMenu(WindowWidth, WindowHeight, GameTitle)
	DrawWaitingMenu(WindowWidth, WindowHeight)
	DrawGameMenu(WindowWidth, WindowHeight)
}

func DrawGrid(col int) {
	dc := gg.NewContext(gridSize, gridSize)

	// Draw the background
	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Draw the lines
	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	cellSize := float64(gridSize) / float64(col)

	for i := 1; i < col; i++ {
		pos := float64(i) * cellSize

		// Vertical line
		dc.DrawLine(pos, 0, pos, float64(gridSize))

		// Horizontal line
		dc.DrawLine(0, pos, float64(gridSize), pos)
	}
	dc.Stroke()

	// Outer border
	offset := lineWidth / 2
	dc.DrawRectangle(offset, offset, float64(gridSize)-lineWidth, float64(gridSize)-lineWidth)
	dc.Stroke()

	GridImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawCircle() {
	dc := gg.NewContext(gridSize, gridSize)

	centerX := float64(cellSize + cellSize/2)
	centerY := float64(cellSize + cellSize/2)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.DrawCircle(centerX, centerY, symbolLength)

	dc.Stroke()

	CircleImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawCross() {
	dc := gg.NewContext(gridSize, gridSize)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	centerX := float64(cellSize + cellSize/2)
	centerY := float64(cellSize + cellSize/2)

	// Diagonal from top-left to bottom-right \
	dc.DrawLine(centerX-symbolLength, centerY-symbolLength, centerX+symbolLength, centerY+symbolLength)

	// Diagonal from bottom-left to top-right /
	dc.DrawLine(centerX-symbolLength, centerY+symbolLength, centerX+symbolLength, centerY-symbolLength)

	dc.Stroke()

	CrossImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawWaitingWheel() {
	const S = 64
	dc := gg.NewContext(S, S)

	cx, cy := float64(S)/2, float64(S)/2
	radius := 22.0
	count := 12

	for i := 0; i < count; i++ {
		angle := float64(i) * (2 * math.Pi) / float64(count)
		x := cx + math.Cos(angle)*radius
		y := cy + math.Sin(angle)*radius

		progress := float64(i) / float64(count)

		r := 2.0 + (3.0 * progress)

		alpha := uint8(50 + (205 * progress))

		col := color.RGBA{R: 0, G: 0, B: 0, A: alpha}

		dc.SetColor(col)
		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	WheelImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawMainMenu(width, height int, title string) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Load the font
	if err := dc.LoadFontFace("client/assets/font.ttf", 48); err != nil {
		log.Println("warning, couldn't load the font")
	}

	// Game title
	dc.SetHexColor("#2C3E50")
	dc.DrawStringAnchored(title, float64(width/2), float64(height)/5, 0.5, 0.5)

	MainMenuImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawWaitingMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Load the font
	if err := dc.LoadFontFace("client/assets/font.ttf", 20); err != nil {
		log.Println("warning, couldn't load the font")
	}

	dc.SetHexColor("#2C3E50")
	dc.DrawStringAnchored("Waiting for another player...", float64(width/2), float64(height)/5, 0.5, 0.5)

	WaitingMenuImage = ebiten.NewImageFromImage(dc.Image())
}

func DrawGameMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Load the font
	if err := dc.LoadFontFace("client/assets/font.ttf", 20); err != nil {
		log.Println("warning, couldn't load the font")
	}

	dc.SetHexColor("#2C3E50")
	dc.DrawStringAnchored("Playing Goonker", (float64(width/2)-(gridSize/2))/2, float64(height)/5, 0.5, 0.5)

	GameMenuImage = ebiten.NewImageFromImage(dc.Image())
}
