package spaceshark

import (
	"image/color"
	"time"

	"github.com/atolVerderben/tentsuyu"
)

type SpaceScene struct {
	starBackground *backgroundImageParts
	ships          []*SpaceShip
	tick           int
	startTickTime  time.Time
	gameStateMsg   GameStateMsg
}

func NewSpaceScene(g *Game) *SpaceScene {

	gm := &SpaceScene{
		ships: []*SpaceShip{},
	}

	imgCursor, _ := tentsuyu.RenderTextToImage([]string{"`*"}, 26, 26, 26.0, tentsuyu.Components.ReturnFont(FntGoRegular), color.White)
	tentsuyu.Components.UIController.SetCustomCursor(26, 26, 0, 0, imgCursor)
	tentsuyu.Components.UIController.Cursor.SetStyle(tentsuyu.CursorPointer)

	gm.ships = append(gm.ships, g.player.homeShip)
	gm.starBackground = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("stars"), count: 20}
	gm.starBackground.SetSize(3200, 3200)

	/*go func(gm *SpaceScene, p *Player) {
		timer := time.NewTicker(50 * time.Millisecond)
		for now := range timer.C {
			// entity updates
			// this is called every 100 millisecondes
			gm.SendMessage(p, now)
		}
	}(gm, g.player)*/

	return gm
}

func (gm *SpaceScene) Update(g *Game) error {

	if gm.tick == 0 {
		gm.startTickTime = time.Now()
	}
	gm.tick++
	//log.Printf("Client Tick: %v", gm.tick)
	tentsuyu.Components.Camera.FollowPlayer(g.player.homeShip, 3200, 3200)
	g.player.Update(g)
	for _, s := range gm.ships {
		s.Update(gm.tick) //time.Now().UnixNano() - gm.startTickTime.UnixNano())
	}
	tentsuyu.Components.UIController.Update()
	return nil
}

func (gm *SpaceScene) Draw(g *Game) error {
	gm.starBackground.Draw(g.screen, false)
	for _, s := range gm.ships {
		s.Draw(g.screen)
	}
	tentsuyu.Components.UIController.Draw(g.screen)
	//g.player.homeShip.Draw(g.screen)
	return nil
}

func (gm *SpaceScene) Msg() GameStateMsg {
	return gm.gameStateMsg
}

func (gm *SpaceScene) SetMsg(gameStateMsg GameStateMsg) {
	gm.gameStateMsg = gameStateMsg
}
