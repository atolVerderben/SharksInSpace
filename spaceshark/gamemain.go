package spaceshark

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
)

type GameMain struct {
	gameStateMsg GameStateMsg
	world        *World
}

func NewGameMain(game *Game) *GameMain {
	g := &GameMain{}
	imgCursor, _ := renderTextToImage([]string{"`*"}, 26, 26, 26.0, tentsuyu.Components.ReturnFont(FntGoRegular), color.White)
	tentsuyu.Components.UIController.SetCustomCursor(26, 26, 0, 0, imgCursor)
	tentsuyu.Components.UIController.Cursor.SetStyle(tentsuyu.CursorPointer)

	g.world = &World{
		characters: []Character{},
		sharks:     CreateSharkController(),
	}

	return g
}

func (g *GameMain) Update(game *Game) error {

	/*if tentsuyu.Input.Button("Left").Down() {
		game.player.Move(Left)
		//p.Facing = Left
		game.player.moving = true
		game.player.moveDelay = 10
	}
	if tentsuyu.Input.Button("Right").Down() {
		game.player.Move(Right)
		//p.Facing = Right
		game.player.moving = true
		game.player.moveDelay = 10
	}
	if tentsuyu.Input.Button("Down").Down() {
		game.player.Move(Down)
		//p.Facing = Right
		game.player.moving = true
		game.player.moveDelay = 10
	}
	if tentsuyu.Input.Button("Up").Down() {
		game.player.Move(Up)
		//p.Facing = Right
		game.player.moving = true
		game.player.moveDelay = 10
	}
	if tentsuyu.Input.LeftClick().JustPressed() {
		game.player.Fire(g.world)

	}
	game.player.Update(game)*/
	g.world.Update()
	//tentsuyu.Components.Camera.FollowPlayer(game.player, 500, 500)

	tentsuyu.Components.UIController.Update()
	return nil
}

func (g *GameMain) Draw(game *Game) error {
	game.background.Draw(game.screen, true)
	g.world.Draw(game.screen)
	tentsuyu.Components.UIController.Draw(game.screen)
	return nil
}

func (g *GameMain) Msg() GameStateMsg {
	return g.gameStateMsg
}

func (g *GameMain) SetMsg(msg GameStateMsg) {
	g.gameStateMsg = msg
}
