package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadBmFont(t *testing.T) {
	f, _ := os.Open("./data/Basic_Latin.json")

	data, _ := ioutil.ReadAll(f)

	font := ReadBmFont(data)

	if font == nil {
		t.FailNow()
	}
}
