package spaceshark

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/atolVerderben/tentsuyu"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"

	assets "github.com/atolVerderben/sharksinspace/spaceshark/internal"
)

var (
	//ZoomLevel is the overall zoom of the game for ease of use
	ZoomLevel        float64
	GameSheet        *SpriteSheet
	PlaySoundEffects = true
	rightAngle       = 1.5708
)

const (
	frameBake int = iota
	frameBomb
	frameGhostGirlSad
	frameGhostGirl
	frameGridBig
	frameGridSmall
	frameKaboom1
	frameKaboom2
	frameKaboom3
	frameKraken
	frameMishipeshu
	frameNureLost
	frameNureWin
	frameNureOnna
	frameBackground
	frameSpaloosh1
	frameSpaloosh2
	frameSpaloosh3
	frameSpalooshLogo
	frameUmibozu
	frameVampHunterSad
	frameVampHunter
	frameVampSad
	frameVamp
)

//Game represents the game itself
type Game struct {
	imageLoadedCh    chan error
	audioLoadedCh    chan error
	gameState        GameState
	pausedState      GameState
	gameData         *GameData
	img              map[string]*ebiten.Image
	orientation      int
	screen           *ebiten.Image
	player           *Player
	background       *backgroundImageParts
	mainState        GameState
	toggleScreenText *tentsuyu.TextElement
}

//NewGame begins a new game of spaloosh!
func NewGame(w, h int) (game *Game, err error) {
	rand.Seed(time.Now().UnixNano())
	defer func() {
		if ferr := finalizeAudio(); ferr != nil && err == nil {
			err = ferr
		}
		if err != nil {
			game = nil
		}
	}()
	game = &Game{
		img:           map[string]*ebiten.Image{},
		imageLoadedCh: make(chan error),
		audioLoadedCh: make(chan error),
		gameData:      NewGameData(GameModeNormalTimed),
	}
	ebiten.SetRunnableInBackground(true) //The timer gets messed up if this isn't enabled

	tentsuyu.BootUp(float64(w), float64(h))

	//GameWorldManager = CreateWorldManager()
	go func() {
		if err := game.loadImages(); err != nil {
			game.imageLoadedCh <- err
		}
		close(game.imageLoadedCh)
	}()
	go func() {
		if err := loadAudio(); err != nil {
			game.audioLoadedCh <- err
		}
		close(game.audioLoadedCh)
	}()

	tentsuyu.Input.RegisterButton("Forward", ebiten.KeyW, ebiten.KeyUp)
	tentsuyu.Input.RegisterButton("Up", ebiten.KeyW, ebiten.KeyUp)
	tentsuyu.Input.RegisterButton("Down", ebiten.KeyS, ebiten.KeyDown)
	tentsuyu.Input.RegisterButton("Backward", ebiten.KeyS, ebiten.KeyDown)
	tentsuyu.Input.RegisterButton("Left", ebiten.KeyA, ebiten.KeyLeft)
	tentsuyu.Input.RegisterButton("Right", ebiten.KeyD, ebiten.KeyRight)
	tentsuyu.Input.RegisterButton("Escape", ebiten.KeyEscape)
	tentsuyu.Input.RegisterButton("ToggleFullscreen", ebiten.KeyF11, ebiten.KeyF)
	tentsuyu.Input.RegisterButton("ChangeScreenScale", ebiten.KeyF10)
	tentsuyu.Input.RegisterButton("ToggleSound", ebiten.KeyF9, ebiten.KeyM)
	tentsuyu.Input.RegisterButton("RotateRight", ebiten.KeyRight)
	tentsuyu.Input.RegisterButton("RotateLeft", ebiten.KeyLeft)
	tentsuyu.Input.RegisterButton("Shoot", ebiten.KeySpace, ebiten.KeyEnter, ebiten.KeyControl)
	tentsuyu.Input.RegisterButton("ZoomIn", ebiten.KeyQ)
	tentsuyu.Input.RegisterButton("ZoomOut", ebiten.KeyE)
	//ZoomLevel = 3.0
	game.player = NewPlayer()

	return
}

func (g *Game) DrawBackground() error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
	/*Sx:     GameSheet.Frames[frameBackground].Frame["x"],
	Sy:     GameSheet.Frames[frameBackground].Frame["y"],
	Width:  GameSheet.Frames[frameBackground].Frame["w"],
	Height: GameSheet.Frames[frameBackground].Frame["h"],*/
	}

	if err := g.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

func (g *Game) loadImages() error {

	font, _ := truetype.Parse(goregular.TTF)
	tentsuyu.Components.UIController.AddFontFile(FntGoRegular, font)
	mono, _ := truetype.Parse(gomono.TTF)
	tentsuyu.Components.UIController.AddFontFile(FntGoMono, mono)
	//tentsuyu.Components.UIController.AddFontFile("font1", tt)

	allImages, err := assets.LoadAllImages()
	if err != nil {
		log.Fatal(err)
	}

	tentsuyu.ImageManager.AddImage("ship", allImages[13])
	tentsuyu.ImageManager.AddImage("shipTurnRight", allImages[12])
	tentsuyu.ImageManager.AddImage("thrust1", allImages[14])
	tentsuyu.ImageManager.AddImage("thrust2", allImages[15])
	tentsuyu.ImageManager.AddImage("stars", allImages[16])
	//tentsuyu.ImageManager.AddImage("background", "assets/sprites/backdrop-purple.png")
	tentsuyu.ImageManager.AddImage("laser1", allImages[6])
	tentsuyu.ImageManager.AddImage("laser2", allImages[7])
	tentsuyu.ImageManager.AddImage("laser3", allImages[8])
	tentsuyu.ImageManager.AddImage("shark1", allImages[9])
	tentsuyu.ImageManager.AddImage("shark2", allImages[10])
	tentsuyu.ImageManager.AddImage("shark3", allImages[11])
	tentsuyu.ImageManager.AddImage("explosion1", allImages[2])
	tentsuyu.ImageManager.AddImage("explosion2", allImages[3])
	tentsuyu.ImageManager.AddImage("explosion3", allImages[4])
	tentsuyu.ImageManager.AddImage("debris1", allImages[0])
	tentsuyu.ImageManager.AddImage("debris2", allImages[1])
	tentsuyu.ImageManager.AddImage("gopher", allImages[5])

	drawImage, err := ebiten.NewImage(64, 64, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	drawImage.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	tentsuyu.ImageManager.AddImage("white", drawImage)
	g.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("white"), count: 20}
	g.background.SetSize(1600, 1600)

	textImg, err := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	textImg.Fill(color.RGBA{R: 238, G: 252, B: 255, A: 255})
	tentsuyu.ImageManager.AddImage("textBubble", textImg)

	healthImg, err := ebiten.NewImage(2, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	healthImg.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	tentsuyu.ImageManager.AddImage("health", healthImg)

	//GameSheet = ReadSpriteSheet("assets/spaloosh-sheet.json")

	//ebiten.SetFullscreen(true)
	return nil

}

//ToggleFullscreen toggles the game in or out of full screen
func (g *Game) ToggleFullscreen() {
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	} else {
		ebiten.SetFullscreen(true)
	}
}

func (g *Game) ToggleSound() {
	PlaySoundEffects = !PlaySoundEffects
}

//Loop controls everything that goes on in the game
func (g *Game) Loop(screen *ebiten.Image) error {
	if g.imageLoadedCh != nil || g.audioLoadedCh != nil {
		select {
		case err := <-g.imageLoadedCh:
			if err != nil {
				return err
			}
			g.imageLoadedCh = nil
		case err := <-g.audioLoadedCh:
			if err != nil {
				return err
			}
			g.audioLoadedCh = nil
		default:
		}
	}
	if g.imageLoadedCh != nil || g.audioLoadedCh != nil {
		return ebitenutil.DebugPrint(screen, "Now Loading...")
	}

	/*if err := audioContext.Update(); err != nil {
		return err
	}*/

	tentsuyu.Input.Update()
	g.screen = screen

	if g.gameState == nil {
		//g.mainState = NewGameMain(g)
		g.gameState = NewTitleMain() //g.mainState
	} else {
		switch g.gameState.Msg() {
		case GameStateMsgReqTitle:
			g.gameState = NewTitleMain()
		case GameStateMsgReqMainMenu:
			g.gameState = NewMainMenu(g)
		case GameStateMsgReqMain:
			g.mainState = NewSpaceScene(g) //NewGameMain(g)
			g.gameState = g.mainState
		case GameStateMsgPause:
			g.pausedState = g.gameState
			g.gameState = CreatePaused()
		case GameStateGameOver:
			g.pausedState = g.gameState
			g.gameState = CreateGameOver()
		case GameStateMsgUnPause:
			g.gameState = g.pausedState
			g.gameState.SetMsg(GameStateMsgNone)
		}
	}
	if g.gameState != nil {
		g.gameState.Update(g)
	}
	tentsuyu.Components.UIController.Update()
	if !ebiten.IsRunningSlowly() {
		if err := g.gameState.Draw(g); err != nil {
			return err
		}
		tentsuyu.Components.UIController.Draw(g.screen)
		//g.toggleScreenText.Draw(g.screen)
	}
	g.gameData.Update()

	if tentsuyu.Input.Button("ToggleFullscreen").JustPressed() {
		g.ToggleFullscreen()
	}
	if tentsuyu.Input.Button("ToggleSound").JustPressed() {
		g.ToggleSound()
	}
	if tentsuyu.Input.Button("ChangeScreenScale").JustPressed() {
		switch g.gameData.currentScale {
		case 1:
			g.gameData.currentScale = 2
			ebiten.SetScreenScale(2)
		case 2:
			g.gameData.currentScale = 1
			ebiten.SetScreenScale(1)
		case 3:
			g.gameData.currentScale = 1
			ebiten.SetScreenScale(1)
		}
	}

	return nil
	//return ebitenutil.DebugPrint(screen, fmt.Sprintf("\nFPS: %.2f", ebiten.CurrentFPS()))
}

//All possible fonts
const (
	FntGoRegular    string = "goregular"
	FntGoBold              = "gobold"
	FntGoItalic            = "goitalic"
	FntGoMono              = "gomono"
	FntGoBoldItalic        = "gobolditalic"
)
