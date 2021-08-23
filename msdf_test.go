package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMSDF(t *testing.T) {
	f, _ := os.Open("./fonts/SignTextNarrow_Bold.ttf")

	data, _ := ioutil.ReadAll(f)

	h := NewFontHandle(data, 8)

	name := h.GetFontName()

	if name == "" {
		t.FailNow()
	}
}
