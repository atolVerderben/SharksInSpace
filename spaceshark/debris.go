package spaceshark

import (
	"math"
	"math/rand"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//debrisController keeps track of all active debriss
type DebrisController struct {
	debris    []*Debris
	maxCount  int
	nextSpawn int
}

//CreateDebrisController creates a new empty debrisController
func CreateDebrisController() *DebrisController {
	bc := &DebrisController{
		debris:   []*Debris{},
		maxCount: 50,
	}
	return bc
}

//AddDebris to the DebrisController
func (bc *DebrisController) AddDebris(debris *Debris) {
	bc.debris = append(bc.debris, debris)
}

//Update all active debriss
func (bc *DebrisController) Update(g *World) {
	if bc.nextSpawn > 0 {
		bc.nextSpawn--
	}
	removeList := []*Debris{}
	for i := range bc.debris {
		remove := bc.debris[i].Update(g)
		if remove {
			removeList = append(removeList, bc.debris[i])
		}
	}

	for r := range removeList {
		g.explosions.CreateExplosion(removeList[r].X, removeList[r].Y, 8, explosionLarge)
		bc.RemoveDebris(removeList[r])
	}
	/*if bc.maxCount > len(bc.debris) && bc.nextSpawn <= 0 {
		log.Println("Added debris")
		rX := rand.Intn(3200)
		rY := rand.Intn(3200)

		for tentsuyu.OnScreen(float64(rX), float64(rY), 109, 109) == true {
			rX = rand.Intn(3200)
			rY = rand.Intn(3200)
		}

		tx, ty := g.player.homeShip.GetPosition()
		angle := math.Atan2(ty-float64(rY), tx-float64(rX))

		bc.AddDebris(CreateDebris(float64(rX), float64(rY), angle))

		r := rand.Intn(100)
		bc.nextSpawn = r + 1

	}*/
}

//Draw all active debriss
func (bc *DebrisController) Draw(screen *ebiten.Image) error {
	for i := range bc.debris {
		x, y := bc.debris[i].GetPosition()
		if tentsuyu.OnScreen(x, y, bc.debris[i].Width, bc.debris[i].Height) {
			if err := bc.debris[i].Draw(screen); err != nil {
				return err
			}
		}
	}
	return nil
}

//Debris is fired from the Tank
type Debris struct {
	desX, desY float64
	ownerID    string
	removeFlag bool
	animCount  int
	ticksAlive int
	health     int
	currImage  string
	*tentsuyu.BasicObject
	imgParts *tentsuyu.BasicImageParts
	rotate   float64
}

//CreateDebris returns a debris object
func CreateDebris(x, y, angle float64) *Debris {
	debris := &Debris{
		health: 1,
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Speed:  float64(rand.Intn(4) + 1),
			Width:  32,
			Height: 32,
			Angle:  angle,
		},
		imgParts: &tentsuyu.BasicImageParts{
			Width:      1,
			Height:     1,
			DestHeight: 1,
			DestWidth:  1,
		},
	}
	debris.VX, debris.VY = debris.GetSpeed()*math.Cos(debris.Angle), debris.GetSpeed()*math.Sin(debris.Angle)
	i := rand.Intn(2)
	if i == 1 {
		debris.currImage = "debris1"
	} else {
		debris.currImage = "debris2"
	}
	return debris
}

func (debris *Debris) ReturnBasicObject() *tentsuyu.BasicObject {
	return debris.BasicObject
}

func (debris *Debris) TakeDamage(dmg int) {
	debris.health -= dmg
}

//Removedebris removes the debris object from the slice
func (bc *DebrisController) RemoveDebris(debris *Debris) {
	//fmt.Println("Attempting to remove debris")
	delete := -1
	for index, e := range bc.debris {
		if e == debris {
			delete = index
			break
		}
	}
	if delete >= 0 {
		bc.debris = append(bc.debris[:delete], bc.debris[delete+1:]...)
	}

}

//InTrajectory runs through the debris's path and sees if it intersects with the passed object
func (debris *Debris) InTrajectory(obj tentsuyu.GameObject) bool {
	vecX, vecY := math.Cos(debris.Angle), math.Sin(debris.Angle)
	x, y := debris.GetPosition()
	for {
		x += debris.GetSpeed() * vecX
		y += debris.GetSpeed() * vecY
		if obj.Contains(x, y) {
			return true
		}
		if x > 1600 || x < 0 || y > 1600 || y < 0 {
			return false
		}
	}
}

//Update the debris -- returns true if debris is removed
func (debris *Debris) Update(g *World) bool {

	if debris.health <= 0 {
		return true
	}
	debris.AddAngle(0.2)
	debris.AddPosition(debris.VX, debris.VY)

	//x, y := debris.GetPosition()
	//log.Printf("debris: (%v,%v)\n", x, y)
	if debris.ticksAlive > 60 {
		for _, char := range g.characters {
			/*if char.IsDead() {
				continue
			}*/

			if tentsuyu.Collision(char.ReturnBasicObject(), debris.BasicObject) {
				char.TakeDamage(1)
				g.explosions.CreateExplosion(debris.X, debris.Y, 8, explosionLarge)
				return true
			}
		}
	}

	if debris.ticksAlive > 60 { // debris may spawn off screen, this is easiest way to have him not immediately disappear
		x, y := debris.GetPosition()
		if x > 3200 {
			debris.SetX(0)
		}
		if x < 0 {
			debris.SetX(3200)

		}
		if y > 3200 {
			debris.SetY(0)
		}
		if y < 0 {
			debris.SetY(3200)
		}
	}

	debris.ticksAlive++
	return false
}

//Draw the debris
func (debris *Debris) Draw(screen *ebiten.Image) error {
	w, h := debris.GetSize()
	op := &ebiten.DrawImageOptions{}
	//op.ImageParts = debris.imgParts
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(debris.Angle)
	op.GeoM.Translate(debris.GetPosition())

	tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(debris.currImage), op); err != nil {
		return err
	}

	//fmt.Println("Drew debris")

	return nil
}
