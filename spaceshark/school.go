package spaceshark

import (
	"math"
	"math/rand"

	"github.com/atolVerderben/tentsuyu/tentsuyutils"
)

type Fish interface {
	GetIDasString() string
	GetX() float64
	GetY() float64
	Update(*World) bool
}

//School is a flock of fish
type School struct {
	fish                    []Fish
	neighborDistanceSquared float64
}

func startNewSchool(x, y float64, num int, randomPlace bool) *School {
	s := &School{
		fish:                    []Fish{},
		neighborDistanceSquared: math.Pow(400, 2),
	}

	for n := 0; n < num; n++ {
		var f Fish
		if randomPlace {
			f = CreateShark(float64(tentsuyutils.RandomBetween(0, 1336)), float64(tentsuyutils.RandomBetween(0, 768)), float64(n)/float64(num)*2.0*math.Pi)
		} else {
			f = CreateShark(x, y, float64(n)/float64(num)*2.0*math.Pi)
		}

		/*angle :=
			f.SetAngle(angle)
		f.VX = f.Speed * math.Cos(angle)
		f.VY = f.Speed * math.Sin(angle)
		if !randomPlace {
			f.X += 40 * f.VX
			f.Y += 40 * f.VY
		}*/

		s.fish = append(s.fish, f)

		//log.Printf("(%v,%v)\n", f.VX, f.VY)
	}

	return s
}

func (s *School) addFish(x, y float64) {

	f := CreateShark(x, y, rand.Float64()*2.0*math.Pi)

	/*f.VX = f.Speed * math.Cos(angle)
	f.VY = f.Speed * math.Sin(angle)*/

	s.fish = append(s.fish, f)
}

//Update the school of fish as it relates to other members of the school
//Call this in conjunction with main gamestate function
func (s *School) Update(predator *SpaceShip) {
	for i := range s.fish {
		var neighbors = []Fish{}
		//  Iterates through all other boids to find neighbors.
		for j := range s.fish {
			if j != i {
				squareDistance := math.Pow(s.fish[j].GetX()-s.fish[i].GetX(), 2) + math.Pow(s.fish[j].GetY()-s.fish[i].GetY(), 2)

				if squareDistance < s.neighborDistanceSquared {
					//log.Printf("Squared: %v < %v\n", squareDistance, s.neighborDistanceSquared)
					neighbors = append(neighbors, s.fish[j])
				}
			}
		}
		//s.fish[i].Update(neighbors, predator)
	}
}

//RemoveFish removes it from the school
func (s *School) RemoveFish(fishID string) {
	//fmt.Println("Attempting to remove Bullet")
	delete := -1
	for index, e := range s.fish {
		if e.GetIDasString() == fishID {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.fish = append(s.fish[:delete], s.fish[delete+1:]...)
	}

}
