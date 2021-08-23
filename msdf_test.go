package fontcatalog

import (
	"image"
	"io/ioutil"
	"os"
	"testing"
)

func TestMSDF(t *testing.T) {
	f, _ := os.Open("./fonts/SignTextNarrow_Bold.ttf")

	data, _ := ioutil.ReadAll(f)

	h := NewFontHandle(data, 32)

	sacle := h.GetScale()

	rgb := image.NewRGBA(image.Rect(0, 0, 32, 32))
	h.GenerateMSDFGlyph(rune('G'), [2]int{32, 32}, rgb.Pix, [2]int{0, 0}, 32*4, [2]float64{0.25, 0.25}, 1, true)

	o, _ := os.Create("./test.png")

	EncodeImage("./png", o, rgb)

	if sacle != 0 {
		t.FailNow()
	}
}
