package spaceshark

type Player struct {
	character Character
	homeShip  *SpaceShip
}

func NewPlayer() *Player {
	p := &Player{
		//character: NewCommonChar(50, 50, FactionPlayer),
		homeShip: NewSpaceShip("Gopher", 30, 50),
	}

	return p
}

func (p *Player) Update(game *Game) {}
