package spaceshark

type Player struct {
	character Character
	homeShip  *SpaceShip
}

func NewPlayer() *Player {
	p := &Player{
		character: NewCommonChar(50, 50, FactionPlayer),
	}

	return p
}

func (p *Player) Update(game *Game) {}
