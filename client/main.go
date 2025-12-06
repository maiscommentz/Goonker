package main

import (
	"Goonker/client/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// main is the entry point of the game
func main() {
	log.Println("Start client")
	game := &Game{}

	ebiten.SetWindowSize(ui.WindowWidth, ui.WindowHeight)
	ebiten.SetWindowTitle(ui.GameTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
