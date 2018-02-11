package spaceshark

import (
	"github.com/atolVerderben/tentsuyu"
)

type Animation struct {
	currFrame, frameCount, frameSpeed int
	imageParts                        *tentsuyu.BasicImageParts
}

func (a *Animation) Update(spriteSheet SpriteSheet) {
	a.frameCount++
	if a.frameCount > a.frameSpeed {
		switch a.currFrame {
		case 1:
			a.currFrame = 2
		case 2:
			a.currFrame = 3
		case 3:
			a.currFrame = 1
		}
		a.frameCount = 0
		a.imageParts.Sx = spriteSheet.Frames[a.currFrame].Frame["x"]
		a.imageParts.Sy = spriteSheet.Frames[a.currFrame].Frame["y"]
	}
}
