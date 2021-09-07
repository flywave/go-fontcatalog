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

	fgeom := NewFontGeometry()

	cs := NewCharsetsASCII()

	fgeom.LoadFromCharset(font, 42, cs)

	glyph := fgeom.GetGlyphFromCodePoint('A')

	index := glyph.GetIndex()
	if index == 0 {
		t.FailNow()
	}

	attr := NewGeneratorAttributes()

	glyph.WrapBox(-1, 2.0, 1.0)
	box := glyph.GetBoxSize()

	if box[0] == 0 {
		t.FailNow()
	}

	bitmap := NewBitmapAlloc(RGBA, [2]int{60, 60})

	err := GlyphGenerator(MOD_MTSDF, bitmap, glyph, attr)

	if err != nil {
		t.FailNow()
	}

	o, err := os.Create("./test.png")

	EncodeImage("./png", o, bitmap.GetImage())

	if err != nil {
		t.FailNow()
	}
}
