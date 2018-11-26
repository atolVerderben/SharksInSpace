package spaceshark

import (
	"math"
	"math/rand"

	"github.com/atolVerderben/flocking/flock/fishtypes"
	"github.com/atolVerderben/tentsuyu"
	"github.com/atolVerderben/tentsuyu/tentsuyutils"
	"github.com/hajimehoshi/ebiten"
)

//BasicSpaceFish is our space fish of the most basic form!
type BasicSpaceFish struct {
	*tentsuyu.BasicObject

	aX, aY             float64 //acceleration
	species            string
	mass               float64
	maxSpeed, maxForce float64
	points             int
	avoiding           bool
	caught             bool
	dead               bool

	hitZone *tentsuyu.BasicObject

	anim *tentsuyu.Animation

	separationDistance float64
	separationStrength float64
	alignmentDistance  float64
	alignmentStrength  float64
	cohesionDistance   float64
	cohesionStrength   float64
}

func newFish(x, y float64, fishType string) *BasicSpaceFish {
	s := &BasicSpaceFish{
		BasicObject:        tentsuyu.NewBasicObject(x, y, 32, 32),
		hitZone:            tentsuyu.NewBasicObject(x, y, 64, 64),
		species:            fishType,
		mass:               1.0,
		separationDistance: 32,
		separationStrength: 2.5,
		cohesionDistance:   128,
		cohesionStrength:   0.75,
		maxSpeed:           3,
		alignmentDistance:  64,
		alignmentStrength:  0.5,
		maxForce:           0.03,
	}

	s.Speed = 2

	switch fishType {
	default:
		s.anim = tentsuyu.NewAnimation(SpaceSheet, []int{FrameShark2, FrameShark3}, 10)
	}

	s.anim.SetCurrentFrame(rand.Intn(2))

	return s
}

func newSalmon(x, y float64) *BasicSpaceFish {
	return newFish(x, y, fishtypes.Salmon)
}

//Update the fish with a list of closest neighbors and a "predator".
//For this game the predator is the boat the player characters are riding in
func (f *BasicSpaceFish) Update(neighbors []*BasicSpaceFish, predator *SpaceShip) {
	if !f.caught {
		f.avoiding = false
		separation := f.calculateSeparation(neighbors)
		separation.Mul(f.separationStrength)
		cohesion := f.calculateCohesion(neighbors)
		cohesion.Mul(f.cohesionStrength)
		alignment := f.calculateAlignment(neighbors)
		alignment.Mul(f.alignmentStrength)
		avoidance := &tentsuyu.Vector2d{}
		if predator != nil {
			avoidance = f.calculatePredator(predator)
		}
		if !f.avoiding {
			f.applyAccelerationForce(separation)
			f.applyAccelerationForce(cohesion)
			f.applyAccelerationForce(alignment)
		} else {
			avoidance.Mul(8.0)
			f.applyAccelerationForce(avoidance)
		}
		f.updatePosition()
		f.anim.Update()

		f.aX = 0
		f.aY = 0
	}
}

func (f *BasicSpaceFish) updatePosition() {
	f.VX += f.aX * f.Speed
	f.VY += f.aY * f.Speed

	angle := math.Atan2(f.VY, f.VX)
	f.SetAngle(angle)

	f.VX = f.Speed * math.Cos(angle)
	f.VY = f.Speed * math.Sin(angle)

	/*if math.Abs(f.VX) > f.maxSpeed {
		if math.Signbit(f.VX) {
			f.VX = -f.maxSpeed
		}
		f.VX = f.maxSpeed

	}

	if math.Abs(f.VY) > f.maxSpeed {
		if math.Signbit(f.VY) {
			f.VY = -f.maxSpeed
		}
		f.VY = f.maxSpeed

	}*/

	f.X += f.VX
	f.Y += f.VY

	f.hitZone.X = f.X
	f.hitZone.Y = f.Y

}

//Draw the FIsh
func (f *BasicSpaceFish) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = f.anim.ReturnImageParts()

	if !f.caught {

		op.ColorM.Scale(0, 0, 0, 1)
	}

	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(-float64(f.anim.ReturnImageParts().Width)/4, -float64(f.anim.ReturnImageParts().Height)/4)
	op.GeoM.Rotate(f.Angle)
	op.GeoM.Translate(f.X, f.Y)

	tentsuyu.ApplyCameraTransform(op, true)

	screen.DrawImage(tentsuyu.ImageManager.ReturnImage(f.species), op)
	return nil
}

//applyAcceleartionForce applies the given force to the acceleration of the fish
func (f *BasicSpaceFish) applyAccelerationForce(force *tentsuyu.Vector2d) {
	//  Acceleration is force devided by mass
	f.aX += force.X / f.mass
	f.aY += force.Y / f.mass
}

func (f *BasicSpaceFish) calculatePredator(predator *SpaceShip) *tentsuyu.Vector2d {
	force := &tentsuyu.Vector2d{}
	//fmt.Printf("Vel: %v \n", predator.velocity.Length())
	/*if predator.velocity.Length() <= 0.03 && !predator.Rotating {
		return force
	}*/
	distance := tentsuyutils.Distance(f.X, f.Y, predator.X, predator.Y)
	if distance <= 150 {
		diff := tentsuyu.VectorSub(f.ReturnVectorPosition(), tentsuyu.Vector2d{X: predator.X, Y: predator.Y})
		//fmt.Print("Before: " + diff.String())
		diff.Normalize()
		//fmt.Print("  Normalized: " + diff.String() + "\n")
		diff.Div(distance)
		force.Add(*diff)
		force.Normalize()
		force.Mul(f.maxSpeed)
		force.Sub(f.ReturnVectorVelocity())
		force.Limit(f.maxForce)
		f.avoiding = true
	}
	return force
}

//calculateSeparation moves the fish apart if they are too close to their neighbors
func (f *BasicSpaceFish) calculateSeparation(neighbors []*BasicSpaceFish) *tentsuyu.Vector2d {
	sep := &tentsuyu.Vector2d{}
	count := 0.0
	for _, n := range neighbors {
		distance := tentsuyutils.Distance(f.X, f.Y, n.X, n.Y)
		if distance < f.separationDistance {
			diff := tentsuyu.VectorSub(f.ReturnVectorPosition(), n.ReturnVectorPosition())
			//fmt.Print("Before: " + diff.String())
			diff.Normalize()
			//fmt.Print("  Normalized: " + diff.String() + "\n")
			diff.Div(distance)
			sep.Add(*diff)
			count++
		}
	}
	if count > 0 {
		sep.Div(count)
	}

	if sep.Length() > 0 {

		//fmt.Println(sep.String())
		//sep.Normalize()
		//sep.Sub(f.ReturnVectorVelocity())
		sep.Normalize()
		sep.Mul(f.maxSpeed)
		sep.Sub(f.ReturnVectorVelocity())
		sep.Limit(f.maxForce)
		//log.Printf("%v\n", sep.ToString())

		//sep.Normalize()
	}

	return sep
}

//seek out the given position. Returns a vector2d of the "steering" force towards the given point
func (f *BasicSpaceFish) seek(target tentsuyu.Vector2d) *tentsuyu.Vector2d {
	desired := tentsuyu.VectorSub(target, f.ReturnVectorPosition())
	desired.Normalize()
	desired.Mul(f.maxSpeed)

	steer := tentsuyu.VectorSub(*desired, f.ReturnVectorVelocity())
	steer.Normalize()
	steer.Limit(f.maxForce)
	return steer
}

//calculateCohesion gives all fish near each other a similar point to follow
func (f *BasicSpaceFish) calculateCohesion(neighbors []*BasicSpaceFish) *tentsuyu.Vector2d {
	cohesion := &tentsuyu.Vector2d{}
	count := 0.0
	for _, n := range neighbors {
		distance := tentsuyutils.Distance(f.X, f.Y, n.X, n.Y)
		if distance < f.cohesionDistance {
			cohesion.Add(n.ReturnVectorPosition())
			count++
		}
	}

	if count > 0 {
		cohesion.Div(count)
		return f.seek(*cohesion)
	}

	return cohesion

}

//calculateAlignment has the fish try to match the velocity of all neighboring fish
func (f *BasicSpaceFish) calculateAlignment(neighbors []*BasicSpaceFish) *tentsuyu.Vector2d {
	align := &tentsuyu.Vector2d{}
	var count float64
	for _, n := range neighbors {
		distance := tentsuyutils.Distance(f.X, f.Y, n.X, n.Y)
		if distance < f.alignmentDistance {
			align.Add(n.ReturnVectorVelocity())
			count++
		}
	}
	if count > 0 {
		align.Div(count)
		align.Normalize()
		align.Mul(f.maxSpeed)
		steer := tentsuyu.VectorSub(*align, f.ReturnVectorVelocity())
		steer.Limit(f.maxForce)
		return steer
	}

	return align
}

//ReturnVectorPosition returns the X,Y position as a vector2d
//This is useful for vector math
func (f BasicSpaceFish) ReturnVectorPosition() tentsuyu.Vector2d {
	return tentsuyu.Vector2d{X: f.X, Y: f.Y}
}

//ReturnVectorVelocity returns the VX,VY as a vector2d
//This is useful for vector math
func (f BasicSpaceFish) ReturnVectorVelocity() tentsuyu.Vector2d {
	return tentsuyu.Vector2d{X: f.VX, Y: f.VY}
}

//HarpoonHit is called when a harpoon collides with a fish
//It takes an angle parameter of the harpoon so the fish can turn perpendicular
func (f *BasicSpaceFish) HarpoonHit(x, y, angle float64) {
	f.Angle = angle + rightAngle
	f.X = x
	f.Y = y
	f.caught = true
}
