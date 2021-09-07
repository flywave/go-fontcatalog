package fontcatalog

import (
	"image"
)

type CharsetData struct {
	font  Charset
	image image.Image
}

type CharsetImage struct {
	data   CharsetData
	width  int
	height int
}

func (i *CharsetImage) Rect() *RectNode {
	return &RectNode{Rect: Rect{0, 0, i.width, i.height}, Index: i.data.font.ID, Rotated: false}
}

func generateImage(render *FontGeometry, char rune, fieldType string, distanceRange int) *CharsetImage {
	var width, height int
	var img image.Image
	var xOffset, yOffset int
	var XAdvance int

	return &CharsetImage{width: width, height: height, data: CharsetData{
		image: img,
		font: Charset{
			ID:       int(char),
			Char:     char,
			Width:    width,
			Height:   height,
			XOffset:  xOffset,
			YOffset:  yOffset,
			XAdvance: XAdvance,
			Channel:  15,
		},
	}}
}
