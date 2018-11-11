package spaceshark

import (
	"math"
	"math/rand"
	"time"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//ShipAction represents the Enum of ship actions
type ShipAction int

//Available actions while in the ship
const (
	ShipStopEngine ShipAction = iota
	ShipStartEngine
	ShipTurnLeft
	ShipStopTurnLeft
	ShipTurnRight
	ShipStopTurnRight
)

//ShipCommand is what will be sent to the server
type ShipCommand struct {
	SeqNum            int        `json:"seqNum"`
	Duration          int64      `json:"duration"`
	InitTime          time.Time  `json:"initTime"`
	Action            ShipAction `json:"action"`
	FromStepBeginning int
}

type ShipType struct {
	X, Y int
}

var (
	ShipShuttle = ShipType{
		X: 28,
		Y: 242,
	}
	ShipTopLeft = ShipType{
		X: 28,
		Y: 42,
	}
)

//SpaceShip represents the spaceship object
type SpaceShip struct {
	ID                                                 string
	Name                                               string
	PrevX, PrevY, AngleOffset, TurnSpeed, Acceleration float64
	Accelerating                                       bool

	BasicImageParts *tentsuyu.BasicImageParts
	CommandQueue    []*ShipCommand
	*tentsuyu.BasicObject
	BoundsWidth, BoundsHeight                          float64
	LastAction                                         ShipAction
	nextCommand                                        *ShipCommand
	IsLocal                                            bool
	PrevAngle                                          float64
	TurningLeft, TurningRight                          bool
	IsServer                                           bool
	health                                             int
	currImage                                          string
	movingAnim                                         int
	coolDown, coolDownMax, hitCoolDown, hitCoolDownMax int
	switchExhaust                                      bool
}

func (ship *SpaceShip) TakeDamage(dmg int) {
	if ship.hitCoolDown == 0 {
		ship.health -= dmg
		ship.hitCoolDown = 1
		tentsuyu.Components.Camera.StartShaking(true)
		AudioPlayer.PlaySE("playerHit")
	}
}

func (ship *SpaceShip) IsDead() bool {
	return ship.health <= 0
}

func ReturnFaction() Faction {
	return FactionPlayer
}

func (ship *SpaceShip) ReturnBasicObject() *tentsuyu.BasicObject {
	return ship.BasicObject
}

func (ship *SpaceShip) SetID(id string) {
	ship.ID = id
}

func (ship SpaceShip) ReturnID() string {
	return ship.ID
}

func (ship *SpaceShip) SetBounds(width, height float64) {
	ship.BoundsWidth = width
	ship.BoundsHeight = height
}

//NewSpaceShip creates the specified SpaceShip
func NewSpaceShip(name string, x, y float64) *SpaceShip {
	s := &SpaceShip{

		Name:           name,
		PrevX:          x,
		PrevY:          y,
		TurnSpeed:      0.07,
		AngleOffset:    -1.5708,
		Acceleration:   .2,
		LastAction:     ShipStopEngine,
		BoundsWidth:    3200,
		BoundsHeight:   3200,
		health:         20,
		coolDownMax:    15,
		hitCoolDownMax: 60,
		/*lead: &Lead{
			BasicImageParts: &tentsuyu.BasicImageParts{
				Sx:     332,
				Sy:     468,
				Width:  32,
				Height: 32,
			},
		},
		Attachments: []*ShipAttachment{},*/

	}

	s.BasicObject = &tentsuyu.BasicObject{
		X:      x,
		Y:      y,
		Width:  64,
		Height: 50,
		Speed:  6.0,
	}

	/*asset := &ShipAttachment{
		angleOffset: AngleDrawOffset * 2,
		BasicImageParts: &tentsuyu.BasicImageParts{
			Width:  21,
			Height: 41,
			Sx:     423,
			Sy:     0,
		},
	}
	s.Attachments = append(s.Attachments, asset)*/

	//First Blue
	/*s.BasicImageParts = &tentsuyu.BasicImageParts{
		Width:      16, //98,
		Height:     16, //75,
		Sx:         shipType.X,
		Sy:         shipType.Y,
		DestHeight: 32,
		DestWidth:  32,
	}*/
	return s
}

//VectorMagnitude is well... |v|
func VectorMagnitude(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

//ExecuteCommand from the CommandQueue
func (ship *SpaceShip) ExecuteCommand(com *ShipCommand) {
	switch com.Action {
	case ShipStartEngine:
		ship.Move()
		ship.Accelerating = true
	case ShipTurnLeft:
		ship.Turn("Left")
		//ship.TurningLeft = true
	case ShipStopTurnLeft:
		ship.TurningLeft = false
	case ShipTurnRight:
		ship.Turn("Right")
		//ship.TurningRight = true
	case ShipStopTurnRight:
		//ship.TurningRight = false
	case ShipStopEngine:
		ship.Accelerating = false
	}

	ship.LastAction = com.Action
}

//Move the Spaceship
func (ship *SpaceShip) Move() {
	ship.VX += ship.Acceleration * math.Cos(ship.Angle)
	ship.VY += ship.Acceleration * math.Sin(ship.Angle)
	maxV := ship.Speed
	if math.Abs(ship.VX) > maxV {
		if ship.VX > 0 {
			ship.VX = maxV
		} else {
			ship.VX = -maxV
		}
	}
	if math.Abs(ship.VY) > maxV {
		if ship.VY > 0 {
			ship.VY = maxV
		} else {
			ship.VY = -maxV
		}
	}
	//ship.VY = Round(ship.VY, 0.05)
	//ship.VX = Round(ship.VX, 0.05)
	ship.Accelerating = true
}

//Turn the ship
func (ship *SpaceShip) Turn(dir string) {
	if dir == "Left" {
		ship.Angle -= ship.TurnSpeed
	}
	if dir == "Right" {
		ship.Angle += ship.TurnSpeed
	}
	ship.Angle = Round(ship.Angle, 0.005)
}

//ShiftCmd removes the first element from the command queue
func (ship *SpaceShip) ShiftCmd() *ShipCommand {
	x, a := ship.CommandQueue[0], ship.CommandQueue[1:]
	ship.CommandQueue = a
	return x
}

//Update the Spaceship
func (ship *SpaceShip) Update() {
	if ship.IsDead() {
		return
	}
	if ship.coolDown > 0 {
		ship.coolDown++
		if ship.coolDown > ship.coolDownMax {
			ship.coolDown = 0
		}
	}

	if ship.hitCoolDown > 0 {
		ship.hitCoolDown++
		if ship.hitCoolDown > ship.hitCoolDownMax {
			ship.hitCoolDown = 0
		}
	}
	//prevX, prevY := ship.GetPosition()
	ship.AddPosition(ship.VX, ship.VY)

	/*if ship.Accelerating {
		ship.Move()
	}
	if ship.TurningLeft {
		ship.Turn("Left")
	}
	if ship.TurningRight {
		ship.Turn("Right")
	}*/

	/*for _, sa := range ship.Attachments {
		sa.Update(prevX, prevY, ship.Angle)
	}*/
	x, y := ship.GetPosition()
	if x > ship.BoundsWidth {
		ship.SetX(0)
	}
	if x < 0 {
		ship.SetX(ship.BoundsWidth)

	}
	if y > ship.BoundsHeight {
		ship.SetY(0)
	}
	if y < 0 {
		ship.SetY(ship.BoundsHeight)
	}
	/*if x > ship.BoundsWidth || x < 0 {
		x = prevX
		ship.SetX(x)
		ship.VX = 0
	}
	if y > ship.BoundsHeight || y < 0 {
		y = prevY
		ship.SetY(y)
		ship.VY = 0
	}*/
	ship.X = Round(ship.X, 0.05)
	ship.Y = Round(ship.Y, 0.05)
	if ship.Accelerating {

		ship.movingAnim++
		if ship.movingAnim > 30 {
			ship.movingAnim = 0
			if rand.Intn(10) >= 3 {
				ship.switchExhaust = true
			} else {
				ship.switchExhaust = false
			}
		}
	} else {
		ship.movingAnim = 0

	}
}

//Draw the Spaceship
func (ship *SpaceShip) Draw(screen *ebiten.Image) error {

	if ship.IsDead() {
		return nil
	}
	if ship.Accelerating {
		op := &ebiten.DrawImageOptions{}
		behindX := 32 * math.Cos(ship.Angle)
		behindY := 32 * math.Sin(ship.Angle)

		op.GeoM.Translate(-float64(64)/2, -float64(24)/2)
		if ship.switchExhaust {
			op.GeoM.Scale(-1, 1)
		}
		op.GeoM.Rotate(ship.Angle + ship.AngleOffset)

		op.GeoM.Translate(ship.X-behindX, ship.Y-behindY)

		if ship.hitCoolDown > 0 {
			op.ColorM.Scale(1, 1, 1, 0.75)
		}
		tentsuyu.ApplyCameraTransform(op, true)
		if ship.movingAnim < 15 {
			if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("thrust1"), op); err != nil {
				return err
			}
		} else {
			if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("thrust2"), op); err != nil {
				return err
			}
		}
	}

	w, h := ship.GetSize()
	op := &ebiten.DrawImageOptions{}
	//op.ImageParts = ship.BasicImageParts
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	if ship.TurningLeft || ship.TurningRight {
		ship.currImage = "shipTurnRight"
		if ship.TurningLeft {
			op.GeoM.Scale(-1, 1)
		}
	} else {
		ship.currImage = "ship"
	}
	op.GeoM.Rotate(ship.Angle + ship.AngleOffset)
	op.GeoM.Translate(ship.GetPosition())
	if ship.hitCoolDown > 0 {
		op.ColorM.Scale(1, 1, 1, 0.5)
	}
	tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(ship.currImage), op); err != nil {
		return err
	}

	//Draw the lead of the Spaceship
	//This is especially useful for UFO type shaped ships to see the front
	/*op = &ebiten.DrawImageOptions{}
	op.ImageParts = ship.lead.BasicImageParts
	op.GeoM.Translate(-float64(ship.lead.BasicImageParts.Width)/2, -float64(ship.lead.BasicImageParts.Height)/2)
	aheadX := 100 * math.Cos(ship.Angle)
	aheadY := 100 * math.Sin(ship.Angle)
	op.GeoM.Translate(aheadX+ship.X, aheadY+ship.Y)
	op.ColorM.Scale(1.0, 1.0, 1.0, 0.4)
	tentsuyu.ApplyCameraTransform(op, true)
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("UI"), op); err != nil {
		return err
	}*/

	return nil
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

func (ship *SpaceShip) ReturnFaction() Faction {
	return FactionPlayer
}
