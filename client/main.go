package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameTitle    = "Goonker"
	WindowWidth  = 600
	WindowHeight = 600
)

func main() {
	log.Println("Start client")
	game := &Game{}
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(GameTitle)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
