package fontcatalog

import (
	"fmt"
	"image"

	"github.com/flywave/imaging"
	"github.com/fogleman/gg"
)

type BitmapFontGenerater struct {
	Opt      BitmapFontOptions
	Charsets []rune
	mapchan  chan []*CharsetImage
	holder   *FontHolder
	font     *FontGeometry
	glyphs   *GlyphGeometryList
}

func NewBitmapFontGenerater(holder *FontHolder, opt BitmapFontOptions) *BitmapFontGenerater {
	glyphs := NewGlyphGeometryList()
	return &BitmapFontGenerater{Opt: opt, mapchan: make(chan []*CharsetImage), holder: holder, glyphs: glyphs, font: NewFontGeometryWithGlyphs(glyphs)}
}

func (g *BitmapFontGenerater) Generate() *BitmapFont {
	font := &BitmapFont{pagesMap: make(map[int]Page), pageSheets: make(map[int]image.Image)}
	start := 0
	done := false
	pageCount := 0
	for done {
		limit := g.Opt.Limit
		if start+g.Opt.Limit > len(g.Charsets) {
			limit = len(g.Charsets) - start
			done = true
		}
		go g.mapCharsets(start, start+limit)
		start += limit
		pageCount++
	}

	for p := 0; p < pageCount; p++ {
		images := <-g.mapchan
		image, chrs := g.packeCharsets(images, p)
		font.Chars = append(font.Chars, chrs...)
		font.pageSheets[p] = image
		var page string
		if pageCount > 1 {
			page = fmt.Sprintf("%s.%d", g.Opt.Filename, p)
		} else {
			page = g.Opt.Filename
		}
		font.Pages = append(font.Pages, page)
	}

	fontmetric := g.holder.getFontInfo()
	fontsize := g.Opt.FontSize
	km := g.font.GetKerning()
	font.Kerning = km.GetKernings()

	charsets := make([]string, len(g.Charsets))
	for i := range g.Charsets {
		charsets[i] = string(g.Charsets[i])
	}

	font.Info = FontInfo{
		Face:         g.font.GetName(),
		Size:         fontsize,
		Bold:         false,
		Italic:       false,
		Charset:      charsets,
		Unicode:      true,
		StretchHeigt: 100,
		Smooth:       1,
		AA:           1,
		Padding:      [4]int{g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.TexturePadding},
		Spacing:      [2]int{g.Opt.FontSpacing[0], g.Opt.FontSpacing[1]},
	}

	font.Common = FontCommon{
		LineHeight:   fontmetric.LineHeight,
		Base:         fontmetric.BaseLine,
		ScaleW:       g.Opt.TextureSize[0],
		ScaleH:       g.Opt.TextureSize[1],
		Pages:        pageCount,
		Packed:       false,
		AlphaChannel: 0,
		RedChannel:   0,
		GreenChannel: 0,
		BlueChannel:  0,
	}

	font.DistanceField = DistanceField{
		FieldType:     g.Opt.FieldType,
		DistanceRange: g.Opt.DistanceRange,
	}

	return font
}

func (g *BitmapFontGenerater) mapCharsets(start, end int) {
	ret := make([]*CharsetImage, g.Opt.Limit)
	for i := start; i < end; i++ {
		ret[i] = generateImage(g.font, g.Charsets[i], g.Opt.FieldType, g.Opt.DistanceRange)
	}
	g.mapchan <- ret
}

func (g *BitmapFontGenerater) packeCharsets(images []*CharsetImage, page int) (image.Image, []Charset) {
	packer := NewMaxRectsBinPacker(g.Opt.TextureSize[0], g.Opt.TextureSize[1], g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.AllowRotation)
	rects := make([]RectNode, len(images))

	for i := range rects {
		rects[i] = *images[i].Rect()
	}
	res := packer.Pack(rects, g.Opt.PackerMethod)

	image := image.NewRGBA(image.Rect(0, 0, res.Width, res.Height))

	dbg := gg.NewContextForImage(image)

	maps := make(map[int]*CharsetImage)
	chars := []Charset{}

	for i := range images {
		maps[images[i].data.font.ID] = images[i]
	}

	for _, node := range res.PlacedRects {
		img := maps[node.Index]
		if node.Rotated {
			img.data.image = imaging.Rotate90(img.data.image)
		}
		fnt := img.data.font
		fnt.X = node.X
		fnt.Y = node.Y
		fnt.Page = page
		chars = append(chars, fnt)
		dbg.DrawImage(img.data.image, node.X, node.Y)
	}

	return dbg.Image(), chars
}
