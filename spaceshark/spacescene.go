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
	world          *World
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
	tentsuyu.Components.Camera.SetZoom(1.0)
	gm.world = &World{
		characters: []Character{},
		sharks:     CreateSharkController(),
		bullets:    CreateBulletController(),
	}
	gm.world.characters = append(gm.world.characters, g.player.homeShip)

	return gm
}

func (gm *SpaceScene) Update(g *Game) error {

	g.player.homeShip.TurningLeft = false
	g.player.homeShip.TurningRight = false
	if tentsuyu.Input.Button("Left").Down() {
		g.player.homeShip.Turn("Left")
		g.player.homeShip.TurningLeft = true
	}
	if tentsuyu.Input.Button("Right").Down() {
		g.player.homeShip.Turn("Right")
		g.player.homeShip.TurningRight = true
	}
	if tentsuyu.Input.Button("Up").Down() {
		g.player.homeShip.Move()
	} else {
		g.player.homeShip.Accelerating = false
	}

	if tentsuyu.Input.Button("Shoot").JustPressed() {
		//gm.world.bullets.AddBullet(CreateBullet("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle))
		gm.world.bullets.Fire("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle)
	}

	if gm.tick == 0 {
		gm.startTickTime = time.Now()
	}
	gm.tick++
	//log.Printf("Client Tick: %v", gm.tick)
	tentsuyu.Components.Camera.FollowPlayer(g.player.homeShip, 3200, 3200)
	g.player.Update(g)
	gm.world.Update()
	/*for _, s := range gm.ships {
		s.Update() //time.Now().UnixNano() - gm.startTickTime.UnixNano())
	}*/
	tentsuyu.Components.UIController.Update()
	return nil
}

func (gm *SpaceScene) Draw(g *Game) error {

	/*op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Width:      256,
		Height:     256,
		DestWidth:  1336,
		DestHeight: 768,
	}
	op.ColorM.Scale(1, 1, 1, 0.5)
	if err := g.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("background"), op); err != nil {
		return err
	}*/
	gm.starBackground.Draw(g.screen, false)
	gm.world.Draw(g.screen)
	/*for _, s := range gm.ships {
		s.Draw(g.screen)
	}*/
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
