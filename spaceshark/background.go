package spaceshark

import "github.com/hajimehoshi/ebiten"
import "github.com/atolVerderben/perihelion/tentsuyu"

type backgroundImageParts struct {
	image         *ebiten.Image
	count         int
	width, height int
}

func (s *backgroundImageParts) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *backgroundImageParts) Len() int {
	w, h := s.image.Size()
	return (s.width/w + 1) * (s.height/h + 2)
}

func (s *backgroundImageParts) Dst(i int) (x0, y0, x1, y1 int) {
	w, h := s.image.Size()
	i, j := i%(s.width/w+1), i/(s.width/w+1)-1
	dx := (-s.count / 4) % w
	dy := (s.count / 4) % h
	dstX := i*w + dx
	dstY := j*h + dy
	return dstX, dstY, dstX + w, dstY + h
}

func (s *backgroundImageParts) Src(i int) (x0, y0, x1, y1 int) {
	w, h := s.image.Size()
	return 0, 0, w, h
}

//Draw the star background
func (s *backgroundImageParts) Draw(screen *ebiten.Image, isStill bool) {
	op := &ebiten.DrawImageOptions{ImageParts: s}
	if !isStill {
		op.ColorM.Scale(1.0, 1.0, 1.0, 0.75)
		tentsuyu.ApplyCameraTransform(op, true)
	}
	screen.DrawImage(s.image, op)
}
