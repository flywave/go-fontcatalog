package fontcatalog

import (
	"bytes"
	"encoding/json"
	"image"
)

type BmCharset struct {
	ID       int     `json:"id"`
	Index    int     `json:"index"`
	Char     rune    `json:"char"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	XOffset  int     `json:"xoffset"`
	YOffset  int     `json:"yoffset"`
	XAdvance int     `json:"xadvance"`
	Channel  Channel `json:"chnl"`
	X        int     `json:"x"`
	Y        int     `json:"y"`
	Page     int     `json:"page"`
}

func (c *BmCharset) Pos() image.Point {
	return image.Pt(c.X, c.Y)
}

func (c *BmCharset) Size() image.Point {
	return image.Pt(c.Width, c.Height)
}

func (c *BmCharset) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: c.Pos(),
		Max: c.Pos().Add(c.Size()),
	}
}

func (c *BmCharset) Offset() image.Point {
	return image.Pt(c.XOffset, c.YOffset)
}

type Channel int

const (
	Blue  Channel = 0x1
	Green Channel = 0x2
	Red   Channel = 0x4
	Alpha Channel = 0x8
	All   Channel = 0xf
)

type CharPair struct {
	First, Second rune
}

type FontInfo struct {
	Face         string   `json:"face"`
	Size         int      `json:"size"`
	Bold         bool     `json:"bold"`
	Italic       bool     `json:"italic"`
	Charset      []string `json:"charset"`
	Unicode      bool     `json:"unicode"`
	StretchHeigt int      `json:"stretchH"`
	Smooth       int      `json:"smooth"`
	AA           int      `json:"aa"`
	Padding      [4]int   `json:"padding"`
	Spacing      [2]int   `json:"spacing"`
}

type FontCommon struct {
	LineHeight   int         `json:"lineHeight"`
	Base         int         `json:"base"`
	ScaleW       int         `json:"scaleW"`
	ScaleH       int         `json:"scaleH"`
	Pages        int         `json:"pages"`
	Packed       bool        `json:"packed"`
	AlphaChannel ChannelInfo `json:"alphaChnl"`
	RedChannel   ChannelInfo `json:"redChnl"`
	GreenChannel ChannelInfo `json:"greenChnl"`
	BlueChannel  ChannelInfo `json:"blueChnl"`
}

func (c *FontCommon) Scale() image.Point {
	return image.Pt(c.ScaleH, c.ScaleH)
}

type ChannelInfo int

const (
	Glyph ChannelInfo = iota
	Outline
	GlyphAndOutline
	Zero
	One
)

type DistanceField struct {
	FieldType     string `json:"fieldType"`
	DistanceRange int    `json:"distanceRange"`
}

type Kerning struct {
	First  rune `json:"first"`
	Second rune `json:"second"`
	Amount int  `json:"amount"`
}

type kerningsort []Kerning

func (k kerningsort) Len() int      { return len(k) }
func (k kerningsort) Swap(i, j int) { k[i], k[j] = k[j], k[i] }
func (k kerningsort) Less(i, j int) bool {
	if k[i].First == k[j].First {
		return k[i].Second < k[j].Second
	}
	return k[i].First < k[j].First
}

type Page struct {
	ID   int
	File string
}

type BmFont struct {
	Pages         []string      `json:"pages"`
	Chars         []BmCharset   `json:"chars"`
	Info          FontInfo      `json:"info"`
	Common        FontCommon    `json:"common"`
	DistanceField DistanceField `json:"distanceField"`
	Kerning       []Kerning     `json:"kernings,omitempty"`

	pagesMap   map[int]Page
	charsMap   map[rune]BmCharset
	kerningMap map[CharPair]int
	pageSheets map[int]image.Image
}

func (ur *BmFont) ToJson() (string, error) {
	b, e := json.Marshal(ur)
	return string(b), e
}

func ReadBmFont(datas []byte) *BmFont {
	ts := &BmFont{}
	data := bytes.NewBuffer(datas)
	json.NewDecoder(data).Decode(&ts)
	return ts
}
