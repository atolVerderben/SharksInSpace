package main

import (
	"github.com/atolVerderben/sharksinspace/spaceshark"
	"github.com/hajimehoshi/ebiten"
)

//Dimensions of the game screen
const (
	ScreenHeight = 1280
	ScreenWidth  = 720
)

func main() {
	game, err := spaceshark.NewGame(ScreenHeight, ScreenWidth)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(game.Loop, ScreenHeight, ScreenWidth, 1.0, "Sharks In Space?!"); err != nil {
		panic(err)
	}
}
