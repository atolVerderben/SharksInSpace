package spaceshark

import (
	"math"
	"math/rand"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type ExplosionController struct {
	explosions []*Explosion
}

func CreateExplosionController() *ExplosionController {
	e := &ExplosionController{
		explosions: []*Explosion{},
	}
	return e
}

//Update all active explosions elements
func (e *ExplosionController) Update() {
	removeList := []*Explosion{}
	for i := range e.explosions {
		remove := e.explosions[i].Update()
		if remove {
			removeList = append(removeList, e.explosions[i])
		}
	}

	for r := range removeList {
		e.RemoveExplosion(removeList[r])
	}
}

//Draw all active explosions elements
func (e *ExplosionController) Draw(screen *ebiten.Image) error {
	for i := range e.explosions {
		if err := e.explosions[i].Draw(screen); err != nil {
			return err
		}
	}
	return nil
}

//Addexplosions to the ExplosionController
func (e *ExplosionController) AddExplosion(explosion *Explosion) {
	e.explosions = append(e.explosions, explosion)
}

//RemoveExplosion removes the bullet object from the slice
func (e *ExplosionController) RemoveExplosion(s *Explosion) {
	delete := -1
	for index, e := range e.explosions {
		if e == s {
			delete = index
			break
		}
	}
	if delete >= 0 {
		e.explosions = append(e.explosions[:delete], e.explosions[delete+1:]...)
	}

}

type Explosion struct {
	owner *tentsuyu.GameObject
	*tentsuyu.BasicObject
	imageParts                 *tentsuyu.BasicImageParts
	dmgAmount                  int
	sx, sy, sx1, sy1, sx2, sy2 int
	w, h, w1, h1, w2, h2       int
	stage, countDown, duration int
	currImage                  string
	size                       explosionSize
}

//Update returns true if the explosion is being removed
func (e *Explosion) Update() bool {
	e.countDown++

	if e.countDown > e.duration {
		if e.size == explosionSmall {
			return true
		}
		if e.stage >= 2 {
			//RemoveSmoke(s)
			return true
		}
		e.stage++

		e.countDown = 0
	}
	return false
}

func (e *Explosion) Draw(screen *ebiten.Image) error {
	w, h := e.GetSize()
	op := &ebiten.DrawImageOptions{}
	//op.ImageParts = e.imageParts
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

	op.GeoM.Rotate(e.Angle)
	op.GeoM.Translate(e.GetPosition())

	//GameCamera.DrawCameraTransform(op)
	tentsuyu.ApplyCameraTransform(op, true)
	if e.stage == 0 {
		e.currImage = "explosion1"
	} else if e.stage == 1 {
		e.currImage = "explosion2"
	} else if e.stage == 2 {
		e.currImage = "explosion3"
	}
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(e.currImage), op); err != nil {
		return err
	}

	return nil
}

//CreateExplosion at the specified x,y coords
func CreateExplosion(x, y float64, duration int, size explosionSize) *Explosion {
	e := &Explosion{
		duration: duration,
		size:     size,
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Width:  64,
			Height: 64,
			Angle:  ((float64(rand.Intn(360)) * math.Pi) / 180.0),
		},
	}

	return e

}

func (e *ExplosionController) CreateExplosion(x, y float64, duration int, size explosionSize) {
	e.AddExplosion(CreateExplosion(x, y, duration, size))
}

type explosionSize int

const (
	explosionSmall explosionSize = iota
	explosionLarge
)
