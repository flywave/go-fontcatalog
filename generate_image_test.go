package fontcatalog

import (
	"os"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	opts := DefaultBitmapFontOptions("")
	render := NewGlyphRenderWithData([]byte(notosans_regular), opts.FontSize)

	charimage := generateImage(render, '�', opts.FontSize, opts.FieldType, opts.DistanceRange)

	if charimage == nil {
		t.FailNow()
	}

	o, _ := os.Create("./test.png")

	EncodeImage("./png", o, charimage.data.image)
}
