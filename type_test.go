package fontcatalog

import (
	"os"
	"testing"
)

func TestUnicodeRanges(t *testing.T) {
	ranges := ReadUnicodeRanges()
	if len(ranges) < 10 {
		t.FailNow()
	}
	ur := ranges[10]

	if ur.Category == "" {
		t.FailNow()
	}
}

func TestReadFontCatalog(t *testing.T) {
	f, _ := os.Open("./data/Default_FontCatalog.json")

	font := ReadFontCatalog(f)

	if font == nil {
		t.FailNow()
	}
}

func TestReadFontCatalogDescription(t *testing.T) {
	f, _ := os.Open("./DefaultFonts.json")

	font := ReadFontCatalogDescription(f)

	if font == nil {
		t.FailNow()
	}
}
