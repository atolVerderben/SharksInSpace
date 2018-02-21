package spaceshark

import (
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten"

	"github.com/atolVerderben/tentsuyu"
)

type TitleMain struct {
	gameStateMsg   GameStateMsg
	timer          int
	offsetX        int
	offsetY        int
	lunkerMode     bool
	lunkerCommand  int
	title          *tentsuyu.TextElement
	menu           *tentsuyu.Menu
	background     *backgroundImageParts
	startButton    *tentsuyu.BasicObject
	currMenu       string
	desc           *tentsuyu.TextElement
	powered        *tentsuyu.TextElement
	start          *tentsuyu.MenuElement
	starBackground *backgroundImageParts
	animCount      int
	controls       *tentsuyu.TextElement
}

func NewTitleMain() *TitleMain {
	imgCursor, _ := tentsuyu.RenderTextToImage([]string{"`*"}, 26, 26, 26.0, tentsuyu.Components.ReturnFont(FntGoRegular), color.White)
	tentsuyu.Components.UIController.SetCustomCursor(26, 26, 0, 0, imgCursor)
	tentsuyu.Components.UIController.Cursor.SetStyle(tentsuyu.CursorPointer)

	t := &TitleMain{
		title: tentsuyu.NewTextElement(500, 25, 500, 50, tentsuyu.Components.ReturnFont(FntGoRegular),
			[]string{"Sharks! In Space?!"}, color.White, 32),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(425, 600, 1300, 400, tentsuyu.Components.ReturnFont(FntGoRegular),
			[]string{"I don't know how sharks survive in space,",
				"        but it's your job to stop them!"}, color.White, 24),
		powered: tentsuyu.NewTextElement(0, 74, 1300, 400, tentsuyu.Components.ReturnFont(FntGoRegular),
			[]string{"Powered by Go (using Ebiten)"}, color.White, 12),

		controls: tentsuyu.NewTextElement(955, 10, 1300, 400, tentsuyu.Components.ReturnFont(FntGoRegular),
			[]string{"Controls:",
				"W / Up:                    Thrust",
				"AD / Left Right:         Rotate Ship",
				"Space / Ctrl / Enter:  Fire",
				"Q / E:                       Zoom In / Out",
				"Escape:                    Pause Menu",
				"F11:                         FullScreen"}, color.White, 20),
	}

	t.start = &tentsuyu.MenuElement{
		UIElement:  tentsuyu.NewTextElementStationary(550, 325, 200, 40, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Start Game"}, color.White, 32),
		Action:     func() { t.gameStateMsg = GameStateMsgReqMain },
		Selectable: true,
	}

	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 150, 50, tentsuyu.Components.ReturnFont(FntGoRegular), []string{"Start Game"}, color.Black, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MainMenu"
			BuildStatsMenu()
			tentsuyu.Components.UIController.ActivateMenu("StatMenu")
			tentsuyu.Components.UIController.DeActivateMenu(prevMenu)
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	t.menu = testMenu
	t.startButton = &tentsuyu.BasicObject{
		X:           518,
		Y:           578,
		Width:       294,
		Height:      68,
		NotCentered: true,
	}
	t.currMenu = "A"
	t.starBackground = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("stars"), count: 20}
	t.starBackground.SetSize(1500, 1500)
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *TitleMain) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.animCount++

	if t.animCount > 40 {
		t.animCount = 0
	}

	t.start.Update()
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if tentsuyu.Input.Button("Escape").JustPressed() {
		os.Exit(0)
	}
	if tentsuyu.Input.Button("Shoot").JustPressed() || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton7) {
		t.gameStateMsg = GameStateMsgReqMain
	}
	return nil
}

func (t *TitleMain) Draw(game *Game) error {
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
	//game.DrawBackground() //background.Draw(game.screen, true)
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 5)
	op.GeoM.Translate(90, 80)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("girl"), op); err != nil {
		return err
	}*/

	/*op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     GameSheet.Frames[frameSpalooshLogo].Frame["x"],
		Sy:     GameSheet.Frames[frameSpalooshLogo].Frame["y"],
		Width:  GameSheet.Frames[frameSpalooshLogo].Frame["w"],
		Height: GameSheet.Frames[frameSpalooshLogo].Frame["h"],
	}
	//op.GeoM.Scale(2, 2)
	//op.GeoM.Translate(250, 20)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     GameSheet.Frames[frameKraken].Frame["x"],
		Sy:     GameSheet.Frames[frameKraken].Frame["y"],
		Width:  GameSheet.Frames[frameKraken].Frame["w"],
		Height: GameSheet.Frames[frameKraken].Frame["h"],
	}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(620, 175)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}*/
	t.starBackground.Draw(game.screen, false)
	t.start.Draw(game.screen)
	t.title.Draw(game.screen)
	t.desc.Draw(game.screen)
	t.controls.Draw(game.screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-32, -25)
	x, y := tentsuyu.Input.GetMouseCoords()

	op.GeoM.Rotate(math.Atan2(y-200, x-650) - rightAngle)
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(650, 200)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("ship"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Rotate(rightAngle * 2)
	op.GeoM.Translate(725, 600)
	currImage := ""
	if t.animCount < 10 {
		currImage = "shark1"
	} else if t.animCount >= 10 && t.animCount < 20 {
		currImage = "shark2"
	} else if t.animCount >= 20 && t.animCount < 30 {
		currImage = "shark1"
	} else if t.animCount >= 30 {
		currImage = "shark3"
	}
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage(currImage), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(2, 2)
	op.GeoM.Translate(35, 10)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("gopher"), op); err != nil {
		return err
	}
	t.powered.Draw(game.screen)

	return nil
}

func (t *TitleMain) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *TitleMain) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
