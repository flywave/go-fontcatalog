package fontcatalog

import (
	"os"
	"testing"
)

func TestFontCatalogGenerater(t *testing.T) {
	data, _ := os.Open("./DefaultFonts.json")
	fcd := ReadFontCatalogDescription(data)

	if fcd == nil {
		t.FailNow()
	}

	opts := DefaultBitmapFontOptions("./test.json")

	gen := NewFontCatalogGenerater(fcd, &opts)

	err := gen.Generate("./data")

	if err != nil {
		t.FailNow()
	}
}
