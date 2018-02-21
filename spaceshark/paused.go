package spaceshark

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Paused is displayed when the game is over
type Paused struct {
	gameStateMsg  GameStateMsg
	timer         int
	offsetX       int
	offsetY       int
	lunkerMode    bool
	lunkerCommand int
	title         *tentsuyu.TextElement
	menu          *tentsuyu.Menu
	background    *backgroundImageParts
	startButton   *tentsuyu.BasicObject
	currMenu      string
	pauseTick     int
	menuTick      int
}

//CreatePaused creates either a winning screen or losing screen depending on the bool value
func CreatePaused() *Paused {
	imgCursor, _ := tentsuyu.RenderTextToImage([]string{"`*"}, 26, 26, 26.0, tentsuyu.Components.ReturnFont(FntGoRegular), color.White)
	tentsuyu.Components.UIController.SetCustomCursor(26, 26, 0, 0, imgCursor)
	tentsuyu.Components.UIController.Cursor.SetStyle(tentsuyu.CursorPointer)
	t := &Paused{
		title:     tentsuyu.NewTextElementStationary(550, 5, 600, 200, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Game Paused"}, color.White, 24),
		pauseTick: 1,
	}

	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 150, 30, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Resume"}, color.White, 22)},
		[]func(){func() { t.gameStateMsg = GameStateMsgUnPause }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 150, 30, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"[R]estart"}, color.White, 22)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMain }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 150, 25, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"[Q]uit Game"}, color.White, 22)},
		[]func(){func() {
			os.Exit(0)
		}})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	//t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("white"), count: 20, width: 1920, height: 1080}
	t.menu = testMenu

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

//Update Paused screen
func (t *Paused) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}

	if t.pauseTick > 0 {
		t.pauseTick++
		if t.pauseTick == 15 {
			t.pauseTick = 0
		}
		return nil
	}

	t.timer++
	if t.menuTick > 0 {
		t.menuTick++
		if t.menuTick > 15 {
			t.menuTick = 0
		}
	}
	if ebiten.GamepadAxis(0, 1) > 0.10 || tentsuyu.Input.Button("Up").JustPressed() {
		if t.menuTick == 0 {
			t.menuTick = 1
			selected, _ := t.menu.ReturnSelected()
			if selected == -1 {
				t.menu.SelectHorizontal(0)
			} else {
				if selected == 0 {
					t.menu.SelectHorizontal(2)
				} else {
					t.menu.SelectHorizontal(selected - 1)
				}
			}
		}
	}
	if ebiten.GamepadAxis(0, 1) < -0.10 || tentsuyu.Input.Button("Down").JustPressed() {
		if t.menuTick == 0 {
			t.menuTick = 1
			selected, _ := t.menu.ReturnSelected()

			if selected == -1 {
				t.menu.SelectHorizontal(0)
			} else {
				if selected == 2 {
					t.menu.SelectHorizontal(0)
				} else {
					t.menu.SelectHorizontal(selected + 1)
				}
			}
		}
	}

	if tentsuyu.Input.Button("Shoot").JustPressed() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) {
		t.menu.PressSelected()
	}

	t.menu.Update()

	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if tentsuyu.Input.Button("Escape").JustPressed() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton7) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) {
		t.gameStateMsg = GameStateMsgUnPause
	}
	if tentsuyu.Input.Button("Restart").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMain
	}
	if tentsuyu.Input.Button("Quit").JustPressed() {
		os.Exit(0)
	}
	return nil
}

//Draw Paused scene
func (t *Paused) Draw(g *Game) error {
	g.pausedState.Draw(g)
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 40)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("map"), op); err != nil {
		return err
	}
	op.GeoM.Translate(0, -40)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("topbar"), op); err != nil {
		return err
	}*/
	/*t.background.Draw(game.screen, false)
	t.menu.Draw(game.screen)
	t.title.Draw(game.screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.3, 0.3)
	op.GeoM.Translate(400, 400)

	tentsuyu.ApplyCameraTransform(op, true)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("shenanijam"), op); err != nil {
		return err
	}*/
	//t.background.Draw(game.screen, false)

	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Width:      32,
		Height:     32,
		DestWidth:  1336,
		DestHeight: 768,
	}
	op.ColorM.Scale(2.11, 2.11, 2.11, 0.25)
	if err := g.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("white"), op); err != nil {
		return err
	}

	t.menu.Draw(g.screen)
	t.title.Draw(g.screen)

	return nil
}

//Msg returns the gamestate msg
func (t *Paused) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *Paused) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
