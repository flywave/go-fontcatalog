package fontcatalog

import (
	"bytes"
	_ "embed"
	"encoding/json"
)

//go:embed unicode-ranges.json
var unicode_ranges string

type UnicodeBlock struct {
	Name  string
	Min   int
	Max   int
	Fonts []string
}

type Font struct {
	Name    string
	Metrics struct {
		Size          int
		DistanceRange int
		Base          int
		LineHeight    int
		LineGap       int
		CapHeight     int
		XHeight       int
	}
	Charset    string
	Bold       *string
	Italic     *string
	BoldItalic *string
}

type FontCatalog struct {
	Name            string
	Type            string
	Size            float64
	MaxWidth        float64
	MaxHeight       float64
	DistanceRange   float64
	Fonts           []Font
	SupportedBlocks []UnicodeBlock
}

var (
	fontCatalog = FontCatalog{
		Name:            "",
		Type:            "",
		Size:            0.0,
		MaxWidth:        0.0,
		MaxHeight:       0.0,
		DistanceRange:   0.0,
		Fonts:           []Font{},
		SupportedBlocks: []UnicodeBlock{},
	}
)

type SdfOptions struct {
	OutputType     string
	Filename       string
	Charset        string
	FontSize       float64
	TexturePadding float64
	FieldType      string
	DistanceRange  float64
	SmartSize      bool
}

var (
	sdfOptions = SdfOptions{
		OutputType:     "json",
		Filename:       "",
		Charset:        "",
		FontSize:       0.0,
		TexturePadding: 2.0,
		FieldType:      "",
		DistanceRange:  0.0,
		SmartSize:      true,
	}
)

type UnicodeRanges struct {
	Category string    `json:"category"`
	Hexrange [2]string `json:"hexrange"`
	Range    [2]int    `json:"range"`
}

func (ur *UnicodeRanges) ToJson() (string, error) {
	b, e := json.Marshal(ur)
	return string(b), e
}

func ReadUnicodeRanges() []UnicodeRanges {
	ts := []UnicodeRanges{}
	data := bytes.NewBuffer([]byte(unicode_ranges))
	json.NewDecoder(data).Decode(&ts)
	return ts
}
