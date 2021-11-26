package fontcatalog

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
)

//go:embed unicode-ranges.json
var unicode_ranges string

type UnicodeBlock struct {
	Name  string   `json:"name"`
	Min   int      `json:"min"`
	Max   int      `json:"max"`
	Fonts []string `json:"fonts"`
}

type FontMetric struct {
	Size          int     `json:"size"`
	DistanceRange float64 `json:"distanceRange"`
	Base          int     `json:"base"`
	LineHeight    int     `json:"lineHeight"`
	LineGap       int     `json:"lineGap"`
	CapHeight     int     `json:"capHeight"`
	XHeight       int     `json:"xHeight"`
}

type Font struct {
	Name       string     `json:"name"`
	Metrics    FontMetric `json:"metrics"`
	Charset    string     `json:"charset"`
	Bold       *string    `json:"blod,omitempty"`
	Italic     *string    `json:"italic,omitempty"`
	BoldItalic *string    `json:"boldItalic,omitempty"`
	Blocks     []string   `json:"blocks,omitempty"`
}

type FontCatalog struct {
	Name            string         `json:"name"`
	Type            string         `json:"type"`
	Size            float64        `json:"size"`
	MaxWidth        float64        `json:"maxWidth"`
	MaxHeight       float64        `json:"maxHeight"`
	DistanceRange   float64        `json:"distanceRange"`
	Fonts           []Font         `json:"fonts"`
	SupportedBlocks []UnicodeBlock `json:"supportedBlocks"`
}

func (ur *FontCatalog) ToJson() (string, error) {
	b, e := json.Marshal(ur)
	return string(b), e
}

func ReadFontCatalog(reader io.Reader) *FontCatalog {
	ts := &FontCatalog{}
	json.NewDecoder(reader).Decode(&ts)
	return ts
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

//go:embed DefaultFonts.json
var default_fonts string

type UnicodeBlockDescription struct {
	Name       string   `json:"name"`
	Bold       *string  `json:"bold,omitempty"`
	Italic     *string  `json:"italic,omitempty"`
	BoldItalic *string  `json:"boldItalic,omitempty"`
	Blocks     []string `json:"blocks"`
}

type FontCatalogDescription struct {
	Name     string                    `json:"name"`
	Size     int                       `json:"size"`
	Distance int                       `json:"distance"`
	Type     string                    `json:"type"`
	FontsDir string                    `json:"fontsDir"`
	Fonts    []UnicodeBlockDescription `json:"fonts"`
}

func (ur *FontCatalogDescription) ToJson() (string, error) {
	b, e := json.Marshal(ur)
	return string(b), e
}

func ReadFontCatalogDescription(reader io.Reader) *FontCatalogDescription {
	ts := &FontCatalogDescription{}
	json.NewDecoder(reader).Decode(&ts)
	return ts
}

var (
	unicodeBlockNames                               = []string{}
	unicodeBlocks                                   = []UnicodeRanges{}
	DefaultFontsDescription *FontCatalogDescription = nil
)

func init() {
	unicodeBlocks = ReadUnicodeRanges()
	for i := range unicodeBlocks {
		unicodeBlockNames = append(unicodeBlockNames, unicodeBlocks[i].Category)
	}
	f := bytes.NewBuffer([]byte(default_fonts))
	DefaultFontsDescription = ReadFontCatalogDescription(f)
}
