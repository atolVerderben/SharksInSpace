package spaceshark

import (
	"math"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//GameBullets holds all active bullets
//var GameBullets = []*Bullet{}

//BulletController keeps track of all active bullets
type BulletController struct {
	bullets []*Bullet
}

//CreateBulletController creates a new empty BulletController
func CreateBulletController() *BulletController {
	bc := &BulletController{
		bullets: []*Bullet{},
	}
	return bc
}

//AddBullet to the BulletController
func (bc *BulletController) AddBullet(bullet *Bullet) {
	bc.bullets = append(bc.bullets, bullet)
}

//Update all active bullets
func (bc *BulletController) Update(g *World) {
	removeList := []*Bullet{}
	for i := range bc.bullets {
		remove := bc.bullets[i].Update(g)
		if remove {
			removeList = append(removeList, bc.bullets[i])
		}
	}

	for r := range removeList {
		g.explosions.CreateExplosion(removeList[r].X, removeList[r].Y, 15, explosionSmall)
		bc.RemoveBullet(removeList[r])
	}
}

//Draw all active bullets
func (bc *BulletController) Draw(screen *ebiten.Image) error {
	for i := range bc.bullets {
		x, y := bc.bullets[i].GetPosition()
		if tentsuyu.OnScreen(x, y, bc.bullets[i].Width, bc.bullets[i].Height) {
			if err := bc.bullets[i].Draw(screen); err != nil {
				return err
			}
		}
	}
	return nil
}

//Bullet is fired from the Tank
type Bullet struct {
	ownerID    string
	removeFlag bool
	animCount  int
	*tentsuyu.BasicObject
}

//CreateBullet returns a bullet object
func CreateBullet(ownerID string, x, y, angle float64) *Bullet {

	bullet := &Bullet{
		ownerID: ownerID,
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Speed:  15,
			Width:  32,
			Height: 75,
			Angle:  angle,
		},
	}

	//log.Printf("Bullet (%v,%v)\n", bullet.desX, bullet.desY)

	//Move in front of tank turrent
	//bullet.AddPosition(20*math.Cos(angle), 20*math.Sin(angle))

	//GameBulletController.AddBullet(bullet)
	//GameBullets = append(GameBullets, bullet)

	//CreateSmoke(bullet.X, bullet.Y, 6, ColorSmokeWhite, true)
	return bullet
}

func (bc *BulletController) Fire(owner string, x, y, angle float64) {
	b1 := CreateBullet(owner, x, y, angle)
	b2 := CreateBullet(owner, x, y, angle)

	adjustX := 24 * math.Cos(angle)
	adjustY := 24 * math.Sin(angle)
	startX := x + adjustY
	startY := y - adjustX
	b1.SetPosition(startX, startY)
	b1.AddPosition((25 * math.Cos(angle)), (25 * math.Sin(angle)))

	adjustX = 12 * math.Cos(angle)
	adjustY = 12 * math.Sin(angle)
	startX = x - adjustY
	startY = y + adjustX
	b2.SetPosition(startX, startY)
	b2.AddPosition((25 * math.Cos(angle)), (25 * math.Sin(angle)))

	bc.AddBullet(b1)
	bc.AddBullet(b2)
}

//RemoveBullet removes the bullet object from the slice
func (bc *BulletController) RemoveBullet(bullet *Bullet) {
	//fmt.Println("Attempting to remove Bullet")
	delete := -1
	for index, e := range bc.bullets {
		if e == bullet {
			delete = index
			break
		}
	}
	if delete >= 0 {
		bc.bullets = append(bc.bullets[:delete], bc.bullets[delete+1:]...)
	}

}

//InTrajectory runs through the bullet's path and sees if it intersects with the passed object
func (bullet *Bullet) InTrajectory(obj tentsuyu.GameObject) bool {
	vecX, vecY := math.Cos(bullet.Angle), math.Sin(bullet.Angle)
	x, y := bullet.GetPosition()
	for {
		x += bullet.GetSpeed() * vecX
		y += bullet.GetSpeed() * vecY
		if obj.Contains(x, y) {
			return true
		}
		if x > 1600 || x < 0 || y > 1600 || y < 0 {
			return false
		}
	}
}

//Update the bullet -- returns true if bullet is removed
func (bullet *Bullet) Update(g *World) bool {

	bullet.AddPosition(bullet.GetSpeed()*math.Cos(bullet.Angle), bullet.GetSpeed()*math.Sin(bullet.Angle))
	bullet.animCount++
	if bullet.animCount >= 45 {
		bullet.animCount = 0
	}
	//x, y := bullet.GetPosition()
	//log.Printf("Bullet: (%v,%v)\n", x, y)

	for _, shark := range g.sharks.sharks {

		/*if char.ReturnFaction() == bullet.owner.ReturnFaction() {
			continue
		}*/

		if tentsuyu.Collision(shark.ReturnBasicObject(), bullet.BasicObject) {
			shark.TakeDamage(1)
			if shark.IsDead() {
				g.player.sharksKilled++
			}
			return true
		}
	}
	for _, d := range g.debris.debris {

		/*if char.ReturnFaction() == bullet.owner.ReturnFaction() {
			continue
		}*/

		if tentsuyu.Collision(d.ReturnBasicObject(), bullet.BasicObject) {
			d.TakeDamage(1)
			return true
		}
	}
	/*for _, guy := range g.Characters {
		char := guy //.ReturnCharacter()
		if tentsuyu.Collision(char.BasicObject, bullet.BasicObject) {
			if char.isDead {
				break
			}
			if bullet.owner == char {
				break

			}
			currTankState := char.isDead
			if bullet.isSmart && char.Contains(bullet.desX, bullet.desY) {
				//CreateSmoke(x, y, 8, ColorSmokeOrange, true)
				if !char.isDead {
					char.TakeDmg(0.5)
					//tank.stats.dmgTaken += 20
					bullet.owner.stats.dmgGiven += 20
					bullet.owner.stats.smartBulletsHitFull++
				}
			} else {
				//CreateSmoke(x, y, 8, ColorSmokeYellow, true)
				if !char.isDead {
					if char.isPlayer {
						char.TakeDmg(0.1)
					} else {
						char.TakeDmg(0.5)
					}
					//tank.stats.dmgTaken += 10
					bullet.owner.stats.dmgGiven += 10
					if bullet.isSmart {
						bullet.owner.stats.smartBulletsHit++
					} else {
						bullet.owner.stats.bulletsHit++
					}
				}
			}
			if currTankState != char.isDead {
				//Just Killed
				bullet.owner.stats.kills++
				g.numEnemies--
			}
			//RemoveBullet(bullet)
			return true
		}
	}*/

	/*for _, ob := range GameObstacles {
		if tentsuyu.Collision(bullet.BasicObject, ob.ReturnBasicObject()) {
			ob.Hit()
			CreateSmoke(bullet.X, bullet.Y, 8, ColorSmokeOrange, true)
			return true
		}
	}*/

	outOfBounds := false
	if bullet.X > 3200 || bullet.X < 0 {
		outOfBounds = true
	}
	if bullet.Y > 3200 || bullet.Y < 0 {
		outOfBounds = true
	}
	if outOfBounds {
		//RemoveBullet(bullet)
		return true
	}
	return false
}

//Draw the bullet
func (bullet *Bullet) Draw(screen *ebiten.Image) error {
	w, h := bullet.GetSize()
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(bullet.Angle - rightAngle)
	op.GeoM.Translate(bullet.GetPosition())
	//if bullet.color == ColorTankBlack {
	op.ColorM.Scale(1, 1, 1, 1)
	//}
	//GameCamera.DrawCameraTransform(op)
	imageName := "white"
	if bullet.animCount < 15 {
		imageName = "laser1"
	} else if bullet.animCount >= 15 && bullet.animCount < 30 {
		imageName = "laser2"
	} else if bullet.animCount >= 30 {
		imageName = "laser3"
	}
	tentsuyu.ApplyCameraTransform(op, true)
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage(imageName), op); err != nil {
		return err
	}

	//fmt.Println("Drew Bullet")

	return nil
}
