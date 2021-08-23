package fontcatalog

import "testing"

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
