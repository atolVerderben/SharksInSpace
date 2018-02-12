package spaceshark

import (
	"math"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type Character interface {
	Draw(*ebiten.Image) error
	Update()
	ReturnBasicObject() *tentsuyu.BasicObject
	TakeDamage(int)
	IsDead() bool
	ReturnFaction() Faction
	SetID(string)
	ReturnID() string
}

type Soldier struct {
	*CommonChar
	weapon Weapon
}

func CreateSoldier(x, y float64, faction Faction) *Soldier {
	s := &Soldier{}
	s.CommonChar = NewCommonChar(x, y, faction)
	s.Speed = 0.6
	s.weapon = CreateBasicGun()
	return s
}

func (s *Soldier) Update() {
	s.moving = false
	s.CommonChar.Update()

	if s.weapon != nil {
		//player.aimX, player.aimY = tentsuyu.Input.GetGameMouseCoords()

		x, y := s.GetPosition()

		if s.aimX-s.X < 0 {
			s.Facing = Left
		} else {
			s.Facing = Right
		}

		s.weapon.SetAngle(math.Atan2(s.aimY-y, s.aimX-x))
	}

	s.weapon.Update()
}

func (s *Soldier) Draw(screen *ebiten.Image) error {
	s.CommonChar.Draw(screen)
	w, h := s.GetSize()
	if s.weapon != nil {
		op := &ebiten.DrawImageOptions{}
		op.ImageParts = s.weapon.ReturnImageParts()
		if s.hitCoolDown > 0 {
			op.ColorM.Scale(1.0, 1.0, 1.0, 0.75)
		}
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		if s.Facing == Right {
			op.GeoM.Rotate(s.weapon.GetAngle())
		}
		if s.Facing == Left {
			op.GeoM.Scale(1, -1)
			op.GeoM.Rotate(s.weapon.GetAngle()) //- c.weapon.AngleOffset())
		}
		//op.GeoM.Rotate(c.weapon.GetAngle())
		op.GeoM.Translate(s.GetPosition())

		tentsuyu.ApplyCameraTransform(op, true)

		if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(s.image), op); err != nil {
			return err
		}
	}
	return nil
}

//Vector2d is a 2D vector
type Vector2d struct {
	X float64
	Y float64
}

//Directions the player can face
var (
	DirNone  = "Forward"
	DirLeft  = "Left"
	DirRight = "Right"
	DirUp    = "Up"
	DirDown  = "Down"
)

//Direction is an int that represents the directions the CommonChar can face
type Direction int

//Available directions
const (
	Up Direction = iota
	Down
	Left
	Right
	Forward
	Backward
	StrafeLeft
	StrafeRight
)

//Move changes the CommonChar's position based on the direction passed
func (c *CommonChar) Move(direction Direction) bool {
	/*if c.isRunning {
		speed += 5
		c.Stamina--
	}*/
	if c.IsDead() {
		return false
	}
	action := true
	c.prevX, c.prevY = c.GetPosition()
	x, y := c.GetPosition()
	switch direction {

	case Up:
		y -= c.Speed
		break

	case Down:
		y += c.Speed
		break

	case Left:
		x -= c.Speed
		break

	case Right:
		x += c.Speed
		break
	}
	c.moving = true
	c.moveDelay = 10

	/*if x > GameWorldmanager.ReturnActiveWorld().WorldWidth || x < 0 {
		x = c.prevX
		action = false
	}
	if y > GameWorldmanager.ReturnActiveWorld().WorldHeight || y < 0 {
		y = c.prevY
		action = false
	}*/

	c.SetPosition(x, y)

	/*for _, tile := range GameWorldmanager.ReturnActiveWorld().EnvCollision {
		//tx := tile.X + float64(tile.BasicObject.Width/2)
		//ty := tile.Y + float64(tile.BasicObject.Height/2)
		if tentsuyu.Collision(tile.BasicObject, c.collisionBox) { //if NearCoords(tx, ty, x, y, float64(tile.BasicObject.Width/2)) {
			x = c.prevX
			y = c.prevY
			action = false
			c.SetPosition(x, y)
			c.collisionBox.X = c.X
			c.collisionBox.Y = c.Y
		}
	}*/

	return action

}

//CommonChar represents the various CommonChars (player and non)
type CommonChar struct {
	*tentsuyu.BasicObject
	Position, Velocity *Vector2d
	Jumping            bool
	imageParts         *tentsuyu.BasicImageParts
	image              string
	hitCoolDown        int
	currWorld          string
	aimX, aimY         float64
	Facing             Direction
	prevX, prevY       float64
	Standing           bool
	Falling            bool
	moving             bool
	moveDelay          int
	anim               *Animation
	destX, destY       float64
	hasDestination     bool
	MaxHealth, Health  int
	faction            Faction
	ID                 string
}

//NewCommonChar returns a pointer to the CommonChar struct at the given location
func NewCommonChar(x, y float64, faction Faction) *CommonChar {
	c := &CommonChar{
		image: "guy",
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Width:  8,
			Height: 8,
			Speed:  0.4,
		},
		Health:  8,
		faction: faction,
		anim: &Animation{
			frameSpeed: 10,
			currFrame:  17,
		},
		imageParts: &tentsuyu.BasicImageParts{
			Sx:         GameSheet.Frames[17].Frame["x"],
			Sy:         GameSheet.Frames[17].Frame["y"],
			Width:      8,
			Height:     8,
			DestWidth:  8,
			DestHeight: 8,
		},
		Facing: Right,
	}
	c.MaxHealth = c.Health
	//c.weapon = CreateBasicGun()
	return c
}

func (c *CommonChar) SetHealth(health int) {
	c.Health = health
	c.MaxHealth = health
}

func (c *CommonChar) SetID(id string) {
	c.ID = id
}

func (c *CommonChar) ReturnID() string {
	return c.ID
}

//Draw the CommonChar
func (c CommonChar) Draw(screen *ebiten.Image) error {

	w, h := c.GetSize()
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = c.imageParts
	if c.hitCoolDown > 0 {
		op.ColorM.Scale(1.0, 1.0, 1.0, 0.75)
	}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	if c.Facing == Left && c.moving {
		op.GeoM.Scale(-1, 1)
	}
	//op.GeoM.Rotate(c.GetAngle() + AngleDrawOffset)
	op.GeoM.Translate(c.GetPosition())
	if c.faction == FactionNeutral {
		op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	}
	if c.faction == FactionHostile {
		op.ColorM.Scale(0.8, 0.4, 0.8, 1)
	}
	tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(c.image), op); err != nil {
		return err
	}

	for i := 0; i < c.Health; i += 2 {
		//Draw Health Bars
		op = &ebiten.DrawImageOptions{}

		//op.GeoM.Translate(-float64(1), -float64(h))

		//op.GeoM.Rotate(c.GetAngle() + AngleDrawOffset)
		mv := 0.0
		switch i {
		case 2:
			mv = 2
		case 4:
			mv = 4
		case 6:
			mv = 6
		}
		op.GeoM.Translate(c.X-float64(c.Width)/2+mv, c.Y+float64(h)/2)
		if c.Health >= 7 {
			op.ColorM.Scale(0, 1, 0, 1)
		} else if c.Health < 7 && c.Health > 2 {
			op.ColorM.Scale(1, 1, 0, 1)
		} else {
			op.ColorM.Scale(1, 0, 0, 1)
		}
		tentsuyu.ApplyCameraTransform(op, true)

		if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("health"), op); err != nil {
			return err
		}
	}

	return nil
}

//Update the CommonChar
func (c *CommonChar) Update() {
	if c.IsDead() {
		return
	}
	if c.hasDestination {
		if c.destX > c.X {
			c.Move(Right)
		} else if c.destX < c.X {
			c.Move(Left)
		}
		if c.destY > c.Y {
			c.Move(Down)
		} else if c.destY < c.Y {
			c.Move(Up)
		}
		if NearCoords(c.X, c.Y, c.destX, c.destY, 1) {
			c.hasDestination = false
			c.X = c.destX
			c.Y = c.destY
		}
	}
	if c.moving {
		c.anim.frameCount++
		//c.imageParts.Sx = GuySheet.Frames[7].Frame["x"]
		//c.imageParts.Sy = GuySheet.Frames[7].Frame["y"]
		if c.anim.frameCount > c.anim.frameSpeed {
			switch c.anim.currFrame {
			case 17:
				c.anim.currFrame = 12
			case 12:
				c.anim.currFrame = 13
			case 13:
				c.anim.currFrame = 12
			}
			c.anim.frameCount = 0
			c.imageParts.Sx = GameSheet.Frames[c.anim.currFrame].Frame["x"]
			c.imageParts.Sy = GameSheet.Frames[c.anim.currFrame].Frame["y"]
		}
	} else {
		if c.moveDelay > 0 {
			c.moveDelay--
		}
		if c.moveDelay == 0 {
			c.anim.currFrame = 17
			c.anim.frameCount = 0
			c.imageParts.Sx = GameSheet.Frames[17].Frame["x"]
			c.imageParts.Sy = GameSheet.Frames[17].Frame["y"]
		}
	}

}

//NearCoords tests if (x,y) are "near" (x1, y1) based on radius
func NearCoords(x, y, x1, y1, radius float64) bool {
	//radius := 20.0
	if x <= (x1+radius) && x >= (x1-radius) {
		if y <= (y1+radius) && y >= (y1-radius) {
			return true
		}
	}

	return false
}

//SetDestination character will try to reach
func (c *CommonChar) SetDestination(x, y float64) {
	c.destX = x
	c.destY = y
	c.hasDestination = true
}

//ReturnFaction of the character
func (c *CommonChar) ReturnFaction() Faction {
	return c.faction
}

//ReturnBasicObject simply returns the basic object of the character
func (c *CommonChar) ReturnBasicObject() *tentsuyu.BasicObject {
	return c.BasicObject
}

func (c *CommonChar) TakeDamage(dmg int) {
	c.Health -= dmg
}

func (c CommonChar) IsDead() bool {
	if c.Health <= 0 {
		return true
	}
	return false
}
