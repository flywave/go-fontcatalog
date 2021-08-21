package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTrueType(t *testing.T) {
	f, _ := os.Open("./fonts/NotoSans-Regular.ttf")

	data, _ := ioutil.ReadAll(f)

	finfo := NewFontInfoFromData(data, 0)

	if finfo == nil {
		t.FailNow()
	}
	size_sdf := float64(26)

	ascent, _, _ := finfo.FontVMetrics()

	scale := finfo.ScaleForPixelHeight(size_sdf)

	baseline := (int)(float64(ascent) * scale)

	if baseline == 0 {
		t.FailNow()
	}
}
