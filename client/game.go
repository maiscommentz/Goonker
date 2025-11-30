package main

import (
	"Goonker/client/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	boardImage *ebiten.Image
)

const (
	// States of the application
	sInit        = 0
	sMenu        = 1
	sGamePlaying = 2
	sGameWin     = 3
	sGameLose    = 4
)

type Game struct {
	menu  *Menu
	state int
}

type Menu struct {
	menuImage *ebiten.Image
	btnPlay   *ui.Button
	btnQuit   *ui.Button
}

/**
 * Update the game state.
 * Called every tick (1/60 [s] by default).
 */
func (g *Game) Update() error {
	switch g.state {
	case sInit:
		g.Init()
	case sMenu:
		if g.menu.btnPlay.IsClicked() {
			g.state = sGamePlaying
		}
		if g.menu.btnQuit.IsClicked() {
			return ebiten.Termination
		}
	}
	return nil
}

/**
 * Draws the game screen.
 * Called every frame (typically 1/60[s] for 60Hz display).
 */
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case sMenu:
		ui.RenderMenu(screen, g.menu.menuImage, g.menu.btnPlay, g.menu.btnQuit)
	case sGamePlaying:
		ui.RenderGame(screen, boardImage)
	case sGameWin:
		ui.RenderGame(screen, boardImage)
	case sGameLose:
		ui.RenderGame(screen, boardImage)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (g *Game) Init() {

	g.menu = &Menu{}

	buttonWidth, buttonHeight := 200.0, 60.0
	centerX := (float64(WindowWidth) - buttonWidth) / 2

	g.menu.btnPlay = ui.NewButton(centerX, 200, buttonWidth, buttonHeight, "Play")
	g.menu.btnQuit = ui.NewButton(centerX, 300, buttonWidth, buttonHeight, "Quit")

	if g.menu.menuImage == nil {
		img := ui.DrawMenu(WindowWidth, WindowHeight, GameTitle)
		g.menu.menuImage = ebiten.NewImageFromImage(img)
	}

	if boardImage == nil {
		grid := ui.DrawGrid()
		boardImage = ebiten.NewImageFromImage(grid)
	}

	g.state = sMenu
}
