package fontcatalog

import (
	"image"
)

type CharsetImage struct {
	font  Charset
	image image.Image
	glyph *GlyphGeometry
}

func generateImage(fgeom *FontGeometry, char rune, fieldType string, distanceRange float64, ec EdgeColoring, angleThreshold float64, seed uint64, attr *GeneratorAttributes) *CharsetImage {
	glyph := fgeom.GetGlyphFromUnicode(char)

	if glyph.IsWhiteSpace() {
		return nil
	}

	glyph.WrapBox(1, distanceRange, 0)

	box := glyph.GetBoxRect()

	xOffset, yOffset, width, height := box[0], box[1], box[2], box[3]

	XAdvance := int(glyph.GetAdvance())

	if fieldType == MOD_MSDF || fieldType == MOD_MTSDF {
		glyph.EdgeColoring(ec, angleThreshold, seed)
	}

	var bitmap *Bitmap
	switch fieldType {
	case MOD_HARD_MASK:
		bitmap = NewBitmapAlloc(GRAY, [2]int{width, height})
	case MOD_SOFT_MASK:
		bitmap = NewBitmapAlloc(GRAY, [2]int{width, height})
	case MOD_SDF:
		bitmap = NewBitmapAlloc(GRAY, [2]int{width, height})
	case MOD_PSDF:
		bitmap = NewBitmapAlloc(GRAY, [2]int{width, height})
	case MOD_MSDF:
		bitmap = NewBitmapAlloc(RGB, [2]int{width, height})
	case MOD_MTSDF:
		bitmap = NewBitmapAlloc(RGBA, [2]int{width, height})
	}

	err := glyphGenerater(fieldType, bitmap, glyph, attr)

	if err != nil {
		return nil
	}

	return &CharsetImage{
		glyph: glyph,
		image: bitmap.GetImage(),
		font: Charset{
			ID:       glyph.GetIndex(),
			Char:     string(char),
			Width:    width,
			Height:   height,
			XOffset:  xOffset,
			YOffset:  yOffset,
			XAdvance: XAdvance,
			Channel:  15,
		},
	}
}
