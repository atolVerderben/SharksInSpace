package main

import (
	"github.com/atolVerderben/sharkstation/spaceshark"
	"github.com/hajimehoshi/ebiten"
)

//Dimensions of the game screen
const (
	ScreenHeight = 1336
	ScreenWidth  = 768
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
