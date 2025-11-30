package ui

import "github.com/hajimehoshi/ebiten/v2"

func RenderMenu(screen *ebiten.Image, menuImage *ebiten.Image, btnPlay, btnQuit *Button) {
	screen.DrawImage(menuImage, nil)
	btnPlay.Draw(screen)
	btnQuit.Draw(screen)
}

func RenderGame(screen *ebiten.Image, boardImage *ebiten.Image) {
	circle := ebiten.NewImageFromImage(DrawCross(0, 0))

	screen.DrawImage(boardImage, nil)
	screen.DrawImage(circle, nil)
}
