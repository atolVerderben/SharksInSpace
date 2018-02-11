package spaceshark

import (
	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Weapon interface to represent any usable weapon for the player and enemies
type Weapon interface {
	Fire() int
	Refill(int)
	ReturnImageParts() *tentsuyu.BasicImageParts
	ReturnImage() *ebiten.Image
	SetAngle(float64)
	GetAngle() float64
	GetWidth() int
	GetHeight() int
	GetX() float64
	GetY() float64
	GetPosition() (float64, float64)
	AngleOffset() float64
	Update()
	SetPosition(float64, float64)
}

//CommonWeapon is the part that follows the mouse and aims
type CommonWeapon struct {
	angle, turnSpeed, x, y float64
	w, h                   int
	imgParts               *tentsuyu.BasicImageParts
	img                    *ebiten.Image
}

func (cw CommonWeapon) GetX() float64 {
	return cw.x
}

func (cw CommonWeapon) GetY() float64 {
	return cw.y
}

func (cw CommonWeapon) GetPosition() (float64, float64) {
	return cw.x, cw.y
}

func (cw CommonWeapon) GetWidth() int {
	return cw.w
}

func (cw CommonWeapon) GetHeight() int {
	return cw.h
}

func (cw *CommonWeapon) GetAngle() float64 {
	return cw.angle
}

func (cw *CommonWeapon) SetAngle(angle float64) {
	cw.angle = angle
}

func (cw *CommonWeapon) ReturnImageParts() *tentsuyu.BasicImageParts {
	return cw.imgParts
}
func (cw *CommonWeapon) ReturnImage() *ebiten.Image {
	return cw.img
}

func (cw *CommonWeapon) SetPosition(x, y float64) {
	cw.x = x
	cw.y = y
}

type BasicGun struct {
	*CommonWeapon
}

func CreateBasicGun() *BasicGun {
	lg := &BasicGun{
		CommonWeapon: &CommonWeapon{
			w: 8,
			h: 8,
			imgParts: &tentsuyu.BasicImageParts{
				Sx:         0,
				Sy:         0,
				Width:      8,
				Height:     8,
				DestWidth:  8,
				DestHeight: 8,
			},
		},
	}
	//lg.img, _ = renderTextToImage([]string{"l"}, lg.w, lg.h, 18, tentsuyu.Components.UIController.ReturnFont("font1"), color.RGBA{R: 35, G: 43, B: 43, A: 255})
	return lg
}

func (lg *BasicGun) AngleOffset() float64 {
	return rightAngle
}

func (lg *BasicGun) Fire() int {
	return 1
}

func (lg *BasicGun) Refill(refill int) {

}

func (lg *BasicGun) Update() {

}
