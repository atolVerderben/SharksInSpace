package assets

import (
	"encoding/base64"
	"image"
	"image/png"

	"github.com/hajimehoshi/ebiten"
)

import "bytes"

//DecodeB64 decodes a base64 string to byte slice
func DecodeB64(message string) []byte {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, _ := base64.StdEncoding.Decode(base64Text, []byte(message))
	//log.Printf("%s\n", base64Text[:l])
	return base64Text[:l]
}

//LoadTankImage loads the tank spritesheet
func LoadTankImage() (*ebiten.Image, error) {

	img, err := png.Decode(bytes.NewReader(DecodeB64("_tanksPNG")))
	if err != nil {
		return nil, err
	}
	img2, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}
	return img2, err
}

//LoadUIImage loads the UI spritesheet
func LoadUIImage() (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(DecodeB64("_uiPNG")))
	if err != nil {
		return nil, err
	}
	img2, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}
	return img2, err
}
