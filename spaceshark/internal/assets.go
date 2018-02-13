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

func LoadImage(imgString string) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(DecodeB64(imgString)))
	if err != nil {
		return nil
	}
	img2, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil
	}
	return img2
}

func LoadAllImages() ([]*ebiten.Image, error) {
	images := []*ebiten.Image{}

	images = append(images, LoadImage(debris))
	images = append(images, LoadImage(debris2))
	images = append(images, LoadImage(explosion1))
	images = append(images, LoadImage(explosion2))
	images = append(images, LoadImage(explosion3))
	images = append(images, LoadImage(gopher))
	images = append(images, LoadImage(laser1))
	images = append(images, LoadImage(laser2))
	images = append(images, LoadImage(laser3))
	images = append(images, LoadImage(shark1))
	images = append(images, LoadImage(shark2))
	images = append(images, LoadImage(shark3))
	images = append(images, LoadImage(shipidle))
	images = append(images, LoadImage(shipidle1))
	images = append(images, LoadImage(thrust1))
	images = append(images, LoadImage(thrust2))
	images = append(images, LoadImage(claystars))

	return images, nil
}
