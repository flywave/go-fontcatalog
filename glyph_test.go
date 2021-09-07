package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGlyphRender(t *testing.T) {
	f, _ := os.Open("./NotoSans-Regular.ttf")

	data, _ := ioutil.ReadAll(f)

	font := NewFontHolder(data)

	if font.m == nil {
		t.FailNow()
	}

	glist := NewGlyphGeometryList()

	fgeom := NewFontGeometryWithGlyphs(glist)

	cs := NewCharsetsASCII()

	n := fgeom.LoadFromCharset(font, 41, cs)

	if n == 0 {
		t.FailNow()
	}

	glyph := fgeom.GetGlyphFromUnicode('J')

	index := glyph.GetIndex()
	if index == 0 {
		t.FailNow()
	}

	attr := NewGeneratorAttributes()

	glyph.WrapBox(1, 4.0, 1.0)
	box := glyph.GetBoxRect()

	glyph.EdgeColoring(EdgeColoringInkTrap, 3.0, 6364136223846793005)

	bitmap := NewBitmapAlloc(RGB, [2]int{box[2], box[3]})

	err := GlyphGenerator(MOD_MSDF, bitmap, glyph, attr)

	if err != nil {
		t.FailNow()
	}

	o, err := os.Create("./test.png")

	EncodeImage("./png", o, bitmap.GetImage())

	if err != nil {
		t.FailNow()
	}
}
