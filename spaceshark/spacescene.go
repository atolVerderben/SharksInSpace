package spaceshark

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"

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
	deadTimer      int
}

func NewSpaceScene(g *Game) *SpaceScene {

	gm := &SpaceScene{
		ships:     []*SpaceShip{},
		hud:       tentsuyu.NewHUD(),
		pauseTick: 1,
	}

	g.player.homeShip.health = 20
	g.player.homeShip.VX = 0
	g.player.homeShip.VY = 0
	g.player.homeShip.hitCoolDown = 0
	g.player.sharksKilled = 0

	g.player.homeShip.SetPosition(float64(rand.Intn(3200)), float64(rand.Intn(3200)))

	tentsuyu.Components.Camera.MaxZoomOut = 0.5

	gm.hud.AddTopLeft(tentsuyu.NewTextElement(0, 0, 200, 30, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Hit Points"}, color.White, 24))
	gm.hud.AddTopLeft(tentsuyu.NewUINumberDisplayInt(&g.player.homeShip.health, 0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoMono), 24, color.White))

	gm.hud.AddTopRight(tentsuyu.NewTextElement(0, 0, 200, 30, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Sharks Destroyed"}, color.White, 24))
	gm.hud.AddTopRight(tentsuyu.NewUINumberDisplayInt(&g.player.sharksKilled, 0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoMono), 24, color.White))

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
	if gm.pauseTick > 0 {
		gm.pauseTick++
		if gm.pauseTick == 15 {
			tentsuyu.Components.UIController.DrawCursor = false
			gm.pauseTick = 0
		}
		return nil
	}

	if tentsuyu.Input.Button("Escape").JustPressed() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton7) {
		if gm.playerJustDied {
			return nil
		}
		gm.gameStateMsg = GameStateMsgPause
		gm.pauseTick = 1
		return nil
	}

	if gm.deadTimer > 0 {
		gm.deadTimer++
		if gm.deadTimer > 60 {
			gm.SetMsg(GameStateGameOver)
			gm.deadTimer = 0
		}
	}

	if g.player.homeShip.IsDead() && !gm.playerJustDied {
		gm.playerJustDied = true
		g.gameData.SetHighScore(g.player.sharksKilled)
		g.gameData.currScore = g.player.sharksKilled
		gm.world.explosions.CreateExplosion(g.player.homeShip.X, g.player.homeShip.Y, 15, explosionLarge)
		gm.deadTimer = 1
	}

	g.player.homeShip.TurningLeft = false
	g.player.homeShip.TurningRight = false
	if tentsuyu.Input.Button("Left").Down() || ebiten.GamepadAxis(0, 2) < -0.40 {
		g.player.homeShip.Turn("Left")
		g.player.homeShip.TurningLeft = true
	}
	if tentsuyu.Input.Button("Right").Down() || ebiten.GamepadAxis(0, 2) > 0.40 {
		g.player.homeShip.Turn("Right")
		g.player.homeShip.TurningRight = true
	}
	if tentsuyu.Input.Button("Up").Down() || ebiten.GamepadAxis(0, 1) > 0.40 {
		g.player.homeShip.Move()
	} else {
		g.player.homeShip.Accelerating = false
	}

	if tentsuyu.Input.Button("Shoot").Down() || ebiten.GamepadAxis(0, 5) > 0.10 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) {
		//gm.world.bullets.AddBullet(CreateBullet("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle))
		if g.player.homeShip.coolDown == 0 {
			g.player.homeShip.coolDown = 1
			gm.world.bullets.Fire("Player", g.player.homeShip.X, g.player.homeShip.Y, g.player.homeShip.Angle)
		}
	}

	if tentsuyu.Input.Button("ZoomOut").Down() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton5) {
		tentsuyu.Components.Camera.ZoomOut()
	}
	if tentsuyu.Input.Button("ZoomIn").Down() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton4) {
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
