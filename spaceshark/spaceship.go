package spaceshark

import (
	"math"
	"time"

	"github.com/atolVerderben/perihelion/tentsuyu"
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
	BoundsWidth, BoundsHeight float64
	LastAction                ShipAction
	nextCommand               *ShipCommand
	IsLocal                   bool
	PrevAngle                 float64
	TurningLeft, TurningRight bool
	IsServer                  bool
	health                    float64
	EngineExhaust             []*Exhaust
}

func (ship *SpaceShip) TakeDamage(dmg float64) {
	ship.health -= dmg
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

//CreateSpaceShip creates the specified SpaceShip
func CreateSpaceShip(name string, x, y float64, shipType ShipType) *SpaceShip {
	s := &SpaceShip{

		Name:         name,
		PrevX:        x,
		PrevY:        y,
		TurnSpeed:    0.07,
		AngleOffset:  0, //1.5708,
		Acceleration: .2,
		LastAction:   ShipStopEngine,
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
		Width:  32,
		Height: 32,
		Speed:  6.0,
	}

	s.EngineExhaust = append(s.EngineExhaust, CreateExhaust("left"), CreateExhaust("right"))

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
	s.BasicImageParts = &tentsuyu.BasicImageParts{
		Width:      16, //98,
		Height:     16, //75,
		Sx:         shipType.X,
		Sy:         shipType.Y,
		DestHeight: 32,
		DestWidth:  32,
	}
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
	ship.VY = Round(ship.VY, 0.05)
	ship.VX = Round(ship.VX, 0.05)
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
func (ship *SpaceShip) Update(tick int) {
	/*if true { //!ship.IsLocal {
		if ship.nextCommand != nil {
			if ship.nextCommand.SeqNum == tick {
				ship.ExecuteCommand(ship.nextCommand)
				ship.nextCommand = nil
			} /*else {
				if ship.nextCommand.SeqNum == 0 {
					ship.nextCommand = nil
				}
			}//--end here
		}
		if ship.nextCommand == nil {
			for len(ship.CommandQueue) > 0 {
				ship.nextCommand = ship.ShiftCmd()
				if ship.nextCommand.SeqNum == tick || ship.nextCommand.SeqNum == 0 {
					if ship.nextCommand.SeqNum == 0 {
						fmt.Println("=========================================== 0! ========================================================")
					}
					ship.ExecuteCommand(ship.nextCommand)
					ship.nextCommand = nil
				} else if ship.nextCommand.SeqNum > 3 {
					fmt.Println("=========================================== 4! ========================================================")
					ship.ExecuteCommand(ship.nextCommand)
					ship.nextCommand = nil
				} else {
					break
				}
				/*if ship.nextCommand.SeqNum == 0 {
					ship.nextCommand.SeqNum = 1
					break
				} else if ship.nextCommand.SeqNum > 3 {
					ship.nextCommand.SeqNum = 3
					break
				} else if ship.nextCommand.SeqNum >= tick {
					ship.ExecuteCommand(ship.nextCommand)
					ship.nextCommand = nil
				} else {
					break
				}//--end here

			}
		}
	} else {
		for len(ship.CommandQueue) > 0 {
			ship.nextCommand = ship.ShiftCmd()
			ship.ExecuteCommand(ship.nextCommand)
			ship.nextCommand = nil

		}
	}*/
	del := []int{}
	for i, c := range ship.CommandQueue {
		if c.SeqNum == tick {
			ship.ExecuteCommand(c)
			del = append(del, i)
			//remove this element
		}
	}
	for x, i := range del {
		i = i - x
		ship.CommandQueue = append(ship.CommandQueue[:i], ship.CommandQueue[i+1:]...)
	}
	prevX, prevY := ship.GetPosition()
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
	if x > ship.BoundsWidth || x < 0 {
		x = prevX
		ship.SetX(x)
		ship.VX = 0
	}
	if y > ship.BoundsHeight || y < 0 {
		y = prevY
		ship.SetY(y)
		ship.VY = 0
	}
	ship.X = Round(ship.X, 0.05)
	ship.Y = Round(ship.Y, 0.05)
	if ship.Accelerating {

		for _, e := range ship.EngineExhaust {
			e.Update(math.Abs(ship.VX), math.Abs(ship.VY), ship.Speed/4)
		}
	}
}

//Draw the Spaceship
func (ship *SpaceShip) Draw(screen *ebiten.Image) error {
	w, h := ship.GetSize()
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = ship.BasicImageParts
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(ship.Angle + ship.AngleOffset)
	op.GeoM.Translate(ship.GetPosition())
	tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("ships"), op); err != nil {
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
	if ship.Accelerating {
		for _, e := range ship.EngineExhaust {
			op = &ebiten.DrawImageOptions{}
			op.ImageParts = e.activeImageParts

			behindX := 16 * math.Cos(ship.Angle)
			behindY := 16 * math.Sin(ship.Angle)
			if e.Name == "right" {

				//Strafe Left
				/*if ebiten.IsKeyPressed(ebiten.KeyA) {
					character.X += character.vY
					character.Y += -character.vX
				}
				//Strafe Right
				if ebiten.IsKeyPressed(ebiten.KeyD) {
					character.X += -character.vY
					character.Y += character.vX


				}*/
				//behindX += -behindY * 20
				//behindY += behindX * 20
			}

			op.GeoM.Translate(-float64(e.activeImageParts.DestWidth)/2, -float64(e.activeImageParts.DestHeight)/2)
			op.GeoM.Rotate(ship.Angle + ship.AngleOffset)
			if e.Name == "left" {
				//behindY = 16 * math.Sin(ship.Angle)
				adjustX := 3 * math.Cos(ship.Angle)
				adjustY := 3 * math.Sin(ship.Angle)
				startX := ship.X + adjustY
				startY := ship.Y - adjustX
				op.GeoM.Translate(startX-behindX, startY-behindY)
			}
			if e.Name == "right" {
				adjustX := 12 * math.Cos(ship.Angle)
				adjustY := 12 * math.Sin(ship.Angle)
				startX := ship.X - adjustY
				startY := ship.Y + adjustX
				op.GeoM.Translate(startX-behindX, startY-behindY)
			}
			tentsuyu.ApplyCameraTransform(op, true)
			if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("ships"), op); err != nil {
				return err
			}
		}
	}
	/*for _, sa := range ship.Attachments {
		op = &ebiten.DrawImageOptions{}
		op.ImageParts = sa.BasicImageParts
		op.GeoM.Translate(-float64(sa.BasicImageParts.Width)/2, -float64(sa.BasicImageParts.Height)/2)
		//aheadX := 100 * math.Cos(sa.Angle)
		//aheadY := 100 * math.Sin(sa.Angle)
		op.GeoM.Translate(sa.X, sa.Y)
		op.GeoM.Rotate(sa.Angle + sa.angleOffset)
		op.ColorM.Scale(1.0, 1.0, 1.0, 0.4)
		tentsuyu.ApplyCameraTransform(op, true)
		if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("SpaceAssets"), op); err != nil {
			return err
		}
	}*/
	return nil
}

type Exhaust struct {
	*tentsuyu.BasicObject
	activeImageParts *tentsuyu.BasicImageParts
	stage            map[int]int
	Name             string
}

func CreateExhaust(side string) *Exhaust {
	e := &Exhaust{
		activeImageParts: &tentsuyu.BasicImageParts{
			Width:      8,
			Height:     8,
			Sx:         64,
			Sy:         492,
			DestHeight: 16,
			DestWidth:  16,
		},
		Name:  side,
		stage: map[int]int{},
	}
	e.stage[0] = 0
	e.stage[4] = 28
	e.stage[3] = 40
	e.stage[2] = 52
	e.stage[1] = 64
	return e
}

func (e *Exhaust) Update(vX, vY, quarterV float64) {
	maxX := 1
	maxY := 1
	if vX > quarterV && vX <= quarterV*2 {
		maxX = 2
	} else if vX > quarterV*2 && vX <= quarterV*3 {
		maxX = 3
	} else if vX >= quarterV*4 {
		maxX = 4
	}

	if vY > quarterV && vY <= quarterV*2 {
		maxY = 2
	} else if vY > quarterV*2 && vY <= quarterV*3 {
		maxY = 3
	} else if vY >= quarterV*4 {
		maxY = 4
	}

	if maxX > maxY {
		e.activeImageParts.Sx = e.stage[maxX]
	} else if maxY > maxX {
		e.activeImageParts.Sx = e.stage[maxY]
	} else {
		e.activeImageParts.Sx = e.stage[maxX]
	}

}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}
