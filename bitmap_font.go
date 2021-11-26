package fontcatalog

import (
	"bytes"
	"encoding/json"
	"image"
)

type Charset struct {
	ID       int     `json:"id"`
	Index    int     `json:"index"`
	Char     string  `json:"char"`
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

func (c *Charset) Pos() image.Point {
	return image.Pt(c.X, c.Y)
}

func (c *Charset) Size() image.Point {
	return image.Pt(c.Width, c.Height)
}

func (c *Charset) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: c.Pos(),
		Max: c.Pos().Add(c.Size()),
	}
}

func (c *Charset) Offset() image.Point {
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
	Packed       int         `json:"packed"`
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
	FieldType     string  `json:"fieldType"`
	DistanceRange float64 `json:"distanceRange"`
}

type Kerning struct {
	First  rune    `json:"first"`
	Second rune    `json:"second"`
	Amount float64 `json:"amount"`
}

type KerningSort []Kerning

func (k KerningSort) Len() int      { return len(k) }
func (k KerningSort) Swap(i, j int) { k[i], k[j] = k[j], k[i] }
func (k KerningSort) Less(i, j int) bool {
	if k[i].First == k[j].First {
		return k[i].Second < k[j].Second
	}
	return k[i].First < k[j].First
}

type Page struct {
	ID   int
	File string
}

type BitmapFont struct {
	Pages         []string            `json:"pages"`
	Chars         []Charset           `json:"chars"`
	Info          FontInfo            `json:"info"`
	Common        FontCommon          `json:"common"`
	DistanceField DistanceField       `json:"distanceField"`
	Kerning       KerningSort         `json:"kernings,omitempty"`
	pagesMap      map[int]Page        `json:"-"`
	pageSheets    map[int]image.Image `json:"-"`
}

func (ur *BitmapFont) ToJson() (string, error) {
	b, e := json.Marshal(ur)
	if true {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return string(b), e
}

func ReadBitmapFont(datas []byte) *BitmapFont {
	ts := &BitmapFont{}
	data := bytes.NewBuffer(datas)
	json.NewDecoder(data).Decode(&ts)
	return ts
}
