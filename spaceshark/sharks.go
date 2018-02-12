package spaceshark

import (
	"log"
	"math"
	"math/rand"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//SharkController keeps track of all active sharks
type SharkController struct {
	sharks    []*Shark
	maxCount  int
	nextSpawn int
}

//CreateSharkController creates a new empty SharkController
func CreateSharkController() *SharkController {
	bc := &SharkController{
		sharks:   []*Shark{},
		maxCount: 50,
	}
	return bc
}

//Addshark to the SharkController
func (bc *SharkController) Addshark(shark *Shark) {
	bc.sharks = append(bc.sharks, shark)
}

//Update all active sharks
func (bc *SharkController) Update(g *World) {
	if bc.nextSpawn > 0 {
		bc.nextSpawn--
	}
	removeList := []*Shark{}
	for i := range bc.sharks {
		remove := bc.sharks[i].Update(g)
		if remove {
			removeList = append(removeList, bc.sharks[i])
		}
	}

	for r := range removeList {
		g.explosions.CreateExplosion(removeList[r].X, removeList[r].Y, 8, explosionLarge)
		bc.RemoveShark(removeList[r])
	}
	if bc.maxCount > len(bc.sharks) && bc.nextSpawn <= 0 {
		log.Println("Added Shark")
		rX := rand.Intn(3200)
		rY := rand.Intn(3200)

		for tentsuyu.OnScreen(float64(rX), float64(rY), 109, 109) == true {
			rX = rand.Intn(3200)
			rY = rand.Intn(3200)
		}

		tx, ty := g.player.homeShip.GetPosition()
		angle := math.Atan2(ty-float64(rY), tx-float64(rX))

		bc.Addshark(CreateShark(float64(rX), float64(rY), angle))

		r := rand.Intn(100)
		bc.nextSpawn = r + 1

	}
}

//Draw all active sharks
func (bc *SharkController) Draw(screen *ebiten.Image) error {
	for i := range bc.sharks {
		x, y := bc.sharks[i].GetPosition()
		if tentsuyu.OnScreen(x, y, bc.sharks[i].Width, bc.sharks[i].Height) {
			if err := bc.sharks[i].Draw(screen); err != nil {
				return err
			}
		}
	}
	return nil
}

//Shark is fired from the Tank
type Shark struct {
	desX, desY float64
	isSmart    bool
	color      int
	ownerID    string
	removeFlag bool
	animCount  int
	ticksAlive int
	health     int
	currImage  string
	*tentsuyu.BasicObject
	imgParts *tentsuyu.BasicImageParts
}

//CreateShark returns a shark object
func CreateShark(x, y, angle float64) *Shark {

	shark := &Shark{
		health: 2,
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Speed:  8,
			Width:  64,
			Height: 109,
			Angle:  angle,
		},
		imgParts: &tentsuyu.BasicImageParts{
			Width:      1,
			Height:     1,
			DestHeight: 1,
			DestWidth:  1,
		},
	}

	return shark
}

func (shark *Shark) ReturnBasicObject() *tentsuyu.BasicObject {
	return shark.BasicObject
}

func (shark *Shark) TakeDamage(dmg int) {
	shark.health -= dmg
}

//RemoveShark removes the shark object from the slice
func (bc *SharkController) RemoveShark(shark *Shark) {
	//fmt.Println("Attempting to remove shark")
	delete := -1
	for index, e := range bc.sharks {
		if e == shark {
			delete = index
			break
		}
	}
	if delete >= 0 {
		bc.sharks = append(bc.sharks[:delete], bc.sharks[delete+1:]...)
	}

}

//InTrajectory runs through the shark's path and sees if it intersects with the passed object
func (shark *Shark) InTrajectory(obj tentsuyu.GameObject) bool {
	vecX, vecY := math.Cos(shark.Angle), math.Sin(shark.Angle)
	x, y := shark.GetPosition()
	for {
		x += shark.GetSpeed() * vecX
		y += shark.GetSpeed() * vecY
		if obj.Contains(x, y) {
			return true
		}
		if x > 1600 || x < 0 || y > 1600 || y < 0 {
			return false
		}
	}
}

//Update the shark -- returns true if shark is removed
func (shark *Shark) Update(g *World) bool {

	if shark.health <= 0 {
		return true
	}
	shark.AddPosition(shark.GetSpeed()*math.Cos(shark.Angle), shark.GetSpeed()*math.Sin(shark.Angle))
	shark.animCount++
	if shark.animCount > 60 {
		shark.animCount = 0
	}
	//x, y := shark.GetPosition()
	//log.Printf("shark: (%v,%v)\n", x, y)

	for _, char := range g.characters {
		/*if char.IsDead() {
			continue
		}*/

		if tentsuyu.Collision(char.ReturnBasicObject(), shark.BasicObject) {
			char.TakeDamage(5)
			g.explosions.CreateExplosion(shark.X, shark.Y, 8, explosionLarge)
			return true
		}
	}

	if shark.ticksAlive > 60 { // Shark may spawn off screen, this is easiest way to have him not immediately disappear
		x, y := shark.GetPosition()
		if x > 3200 {
			shark.SetX(0)
		}
		if x < 0 {
			shark.SetX(3200)

		}
		if y > 3200 {
			shark.SetY(0)
		}
		if y < 0 {
			shark.SetY(3200)
		}
	}

	shark.ticksAlive++
	return false
}

//Draw the shark
func (shark *Shark) Draw(screen *ebiten.Image) error {
	w, h := shark.GetSize()
	op := &ebiten.DrawImageOptions{}
	//op.ImageParts = shark.imgParts
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(shark.Angle - rightAngle)
	op.GeoM.Translate(shark.GetPosition())

	tentsuyu.ApplyCameraTransform(op, true)
	if shark.animCount < 15 {
		shark.currImage = "shark1"
	} else if shark.animCount >= 15 && shark.animCount < 30 {
		shark.currImage = "shark2"
	} else if shark.animCount >= 30 && shark.animCount < 45 {
		shark.currImage = "shark1"
	} else if shark.animCount >= 45 {
		shark.currImage = "shark3"
	}
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(shark.currImage), op); err != nil {
		return err
	}

	//fmt.Println("Drew shark")

	return nil
}
