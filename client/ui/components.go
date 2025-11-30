package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	X, Y, Width, Height float64
	Image               *ebiten.Image
	Text                string
}

func NewButton(x, y, w, h float64, text string) *Button {
	b := &Button{
		X: x, Y: y, Width: w, Height: h,
		Text: text,
	}

	dc := gg.NewContext(int(w), int(h))

	dc.DrawRoundedRectangle(0, 0, w, h, 10)
	dc.SetHexColor("#2C3E50")
	dc.Fill()

	// dc.LoadFontFace("arial.ttf", 24)
	dc.SetHexColor("#FFFFFF")
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, 0.35)

	b.Image = ebiten.NewImageFromImage(dc.Image())

	return b
}

func (b *Button) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Image, opts)
}

func (b *Button) IsClicked() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fmx, fmy := float64(mx), float64(my)

		return fmx >= b.X && fmx <= b.X+b.Width &&
			fmy >= b.Y && fmy <= b.Y+b.Height
	}
	return false
}
