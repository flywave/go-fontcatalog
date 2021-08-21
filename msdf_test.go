package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMSDF(t *testing.T) {
	f, _ := os.Open("./fonts/NotoSans-Regular.ttf")

	data, _ := ioutil.ReadAll(f)

	finfo := NewFontInfoFromData(data, 0)

	if finfo == nil {
		t.FailNow()
	}

	_, bitmap := msdfGlyph(finfo, "A", 32, 32)

	w, _ := os.Create("./test.png")
	EncodeImage("png", w, bitmap)
}
