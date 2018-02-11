package spaceshark

import "github.com/hajimehoshi/ebiten"

type World struct {
	characters    []Character
	sharks        *SharkController
	bullets       *BulletController
	player        *Player
	width, height float64
}

func (w *World) Update() {
	for _, char := range w.characters {
		char.Update()
	}
	w.bullets.Update(w)
	w.sharks.Update(w)
}

func (w *World) Draw(screen *ebiten.Image) error {
	for _, c := range w.characters {
		c.Draw(screen)
	}
	w.bullets.Draw(screen)
	w.sharks.Draw(screen)
	return nil
}
