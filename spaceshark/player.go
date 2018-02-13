package spaceshark

type Player struct {
	character    Character
	homeShip     *SpaceShip
	sharksKilled int
	debrisKilled int
}

func NewPlayer() *Player {
	p := &Player{
		//character: NewCommonChar(50, 50, FactionPlayer),
		homeShip: NewSpaceShip("Gopher", 500, 500),
	}

	return p
}

func (p *Player) Update(game *Game) {

}
