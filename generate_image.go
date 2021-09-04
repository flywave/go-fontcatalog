package fontcatalog

import (
	"image"
	"math"
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

func roundNumber(number float64, digits int) float64 {
	n10 := math.Pow10(digits)
	return math.Trunc((number+0.5/n10)*n10) / n10
}

func generateImage(render *GlyphRender, char rune, fontSize int, fieldType string, distanceRange int) *CharsetImage {
	info := render.GetChar(char)
	scale := fontSize / render.UnitsPerEm
	baseline := render.BaseLine
	pad := distanceRange >> 1
	width := info.Width + pad + pad
	height := info.Height + pad + pad
	xOffset := info.XOffset + pad
	yOffset := info.YOffset + pad

	rgb := image.NewRGBA(image.Rect(0, 0, width, height))

	fh := render.GetFont()
	switch fieldType {
	case MOD_SDF:
		fh.GenerateSDFGlyph(char, [2]int{width, height}, rgb.Pix, [2]float64{float64(xOffset), float64(yOffset)}, float64(distanceRange), false)
	case MOD_PSDF:
		fh.GeneratePSDFGlyph(char, [2]int{width, height}, rgb.Pix, [2]float64{float64(xOffset), float64(yOffset)}, float64(distanceRange), false)
	case MOD_MSDF:
		fh.GenerateMSDFGlyph(char, [2]int{width, height}, rgb.Pix, [2]float64{float64(xOffset), float64(yOffset)}, float64(distanceRange), false)
	}
	return &CharsetImage{width: width, height: height, data: CharsetData{
		image: rgb,
		font: Charset{
			ID:       int(char),
			Char:     char,
			Width:    width,
			Height:   height,
			XOffset:  xOffset - pad,
			YOffset:  yOffset + pad + baseline,
			XAdvance: info.Advance * scale,
			Channel:  15,
		},
	}}
}
