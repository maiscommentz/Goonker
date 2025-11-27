package main

import "github.com/hajimehoshi/ebiten/v2"

func RenderMenu(screen *ebiten.Image, menu *Menu) {
	screen.DrawImage(menu.menuImage, nil)
	menu.btnPlay.Draw(screen)
	menu.btnQuit.Draw(screen)
}

func RenderGame(screen *ebiten.Image, g *Game) {
	circle := ebiten.NewImageFromImage(DrawCross(0, 0))

	screen.DrawImage(boardImage, nil)
	screen.DrawImage(circle, nil)
}
