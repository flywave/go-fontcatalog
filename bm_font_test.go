package fontcatalog

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadBitmapFont(t *testing.T) {
	f, _ := os.Open("./data/Basic_Latin.json")

	data, _ := ioutil.ReadAll(f)

	font := ReadBitmapFont(data)

	if font == nil {
		t.FailNow()
	}
}
