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
	pauseTick      int
	hud            *tentsuyu.HUD
	playerJustDied bool
}

func NewSpaceScene(g *Game) *SpaceScene {

	gm := &SpaceScene{
		ships: []*SpaceShip{},
		hud:   tentsuyu.NewHUD(),
	}

	g.player.homeShip.health = 20
	g.player.homeShip.VX = 0
	g.player.homeShip.VY = 0
	g.player.homeShip.hitCoolDown = 0
	g.player.sharksKilled = 0
	tentsuyu.Components.Camera.MaxZoomOut = 0.5

	gm.hud.AddTopLeft(tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Hit Points"}, color.White, 24))
	gm.hud.AddTopLeft(tentsuyu.NewUINumberDisplayInt(&g.player.homeShip.health, 0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoRegular), 24, color.White))

	gm.hud.AddTopRight(tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Sharks Destroyed"}, color.White, 24))
	gm.hud.AddTopRight(tentsuyu.NewUINumberDisplayInt(&g.player.sharksKilled, 0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoRegular), 24, color.White))

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
		explosions: CreateExplosionController(),
		debris:     CreateDebrisController(),
		player:     g.player,
	}
	gm.world.characters = append(gm.world.characters, g.player.homeShip)
	tentsuyu.Components.UIController.DrawCursor = false
	return gm
}

func (gm *SpaceScene) Update(g *Game) error {
	if tentsuyu.Input.Button("Escape").JustPressed() {
		gm.gameStateMsg = GameStateMsgPause
		gm.pauseTick = 1
		return nil
	}
	if gm.pauseTick > 0 {
		gm.pauseTick++
		if gm.pauseTick == 30 {
			tentsuyu.Components.UIController.DrawCursor = false
			gm.pauseTick = 0
		}
		return nil
	}

	if g.player.homeShip.IsDead() && !gm.playerJustDied {
		gm.playerJustDied = true
		gm.world.explosions.CreateExplosion(g.player.homeShip.X, g.player.homeShip.Y, 15, explosionLarge)
	}

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

	if tentsuyu.Input.Button("Shoot").Down() {
		//gm.world.bullets.AddBullet(CreateBullet("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle))
		if g.player.homeShip.coolDown == 0 {
			g.player.homeShip.coolDown = 1
			gm.world.bullets.Fire("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle)
		}
	}

	if tentsuyu.Input.Button("ZoomOut").Down() {
		tentsuyu.Components.Camera.ZoomOut()
	}
	if tentsuyu.Input.Button("ZoomIn").Down() {
		tentsuyu.Components.Camera.ZoomIn()
	}

	if gm.tick == 0 {
		gm.startTickTime = time.Now()
	}
	gm.tick++
	//log.Printf("Client Tick: %v", gm.tick)
	tentsuyu.Components.Camera.FollowPlayer(g.player.homeShip, 3200, 3200)
	g.player.Update(g)
	gm.world.Update()
	gm.hud.Update()
	/*for _, s := range gm.ships {
		s.Update() //time.Now().UnixNano() - gm.startTickTime.UnixNano())
	}*/

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
	if err := g.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("background"), op); err != nil {
		return err
	}*/
	gm.starBackground.Draw(g.screen, false)
	gm.world.Draw(g.screen)
	/*for _, s := range gm.ships {
		s.Draw(g.screen)
	}*/
	tentsuyu.Components.UIController.Draw(g.screen)
	gm.hud.Draw(g.screen)
	//g.player.homeShip.Draw(g.screen)
	return nil
}

func (gm *SpaceScene) Msg() GameStateMsg {
	return gm.gameStateMsg
}

func (gm *SpaceScene) SetMsg(gameStateMsg GameStateMsg) {
	gm.gameStateMsg = gameStateMsg
}
