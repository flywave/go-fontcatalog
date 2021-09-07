package fontcatalog

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/flywave/imaging"
)

//go:embed NotoSans-Regular.ttf
var notosans_regular string

type FontCatalogGenerater struct {
	opts        *BitmapFontOptions
	fontDesc    *FontCatalogDescription
	fontCatalog *FontCatalog
}

func NewFontCatalogGenerater(desc *FontCatalogDescription, opts *BitmapFontOptions) *FontCatalogGenerater {
	ret := &FontCatalogGenerater{fontDesc: desc, opts: opts, fontCatalog: &FontCatalog{Name: desc.Name, Type: desc.Type, Size: float64(desc.Size), DistanceRange: float64(desc.Distance)}}
	return ret
}

func (g *FontCatalogGenerater) Generate(outputPath string) error {
	sdfOptions := *g.opts
	for _, ufont := range g.fontDesc.Fonts {
		fontPath := path.Join(fmt.Sprintf("%s%s.ttf", g.fontDesc.FontsDir, ufont.Name))
		fontData, err := ioutil.ReadFile(fontPath)
		if err != nil {
			return err
		}
		fontHolder := NewFontHolder(fontData)
		fontInfo := fontHolder.getFontInfo()
		font := &Font{
			Name: "Extra",
			Metrics: FontMetric{
				Size:          sdfOptions.FontSize,
				DistanceRange: sdfOptions.DistanceRange,
				Base:          0.0,
				LineHeight:    0.0,
				LineGap:       int(math.Round(float64(fontInfo.LineGap/fontInfo.UnitsPerEm) * float64(sdfOptions.FontSize))),
				CapHeight: int(math.Round(
					float64(fontInfo.Ascent/fontInfo.UnitsPerEm) * float64(sdfOptions.FontSize))),
				XHeight: 0,
			},
			Charset: "",
		}

		g.createFontAssets(fontHolder, font, g.fontCatalog, fontInfo.CharacterSet, fontPath, false, false, outputPath)

		if ufont.Bold != nil {
			boldFontPath := path.Join(fmt.Sprintf("%s%s.ttf", g.fontDesc.FontsDir, *ufont.Bold))
			boldFontData, err := ioutil.ReadFile(boldFontPath)
			if err != nil {
				return err
			}
			boldFontHolder := NewFontHolder(boldFontData)
			boldFontInfo := boldFontHolder.getFontInfo()
			font.Bold = ufont.Bold
			g.createFontAssets(boldFontHolder, font, g.fontCatalog, boldFontInfo.CharacterSet, boldFontPath, true, false, outputPath)
		}

		if ufont.Italic != nil {
			italicFontPath := path.Join(fmt.Sprintf("%s%s.ttf", g.fontDesc.FontsDir, *ufont.Italic))
			italicFontData, err := ioutil.ReadFile(italicFontPath)
			if err != nil {
				return err
			}
			italicFontHolder := NewFontHolder(italicFontData)
			italicFontInfo := italicFontHolder.getFontInfo()
			font.Italic = ufont.Italic
			g.createFontAssets(italicFontHolder, font, g.fontCatalog, italicFontInfo.CharacterSet, italicFontPath, false, true, outputPath)
		}

		if ufont.BoldItalic != nil {
			boldItalicFontPath := path.Join(fmt.Sprintf("%s%s.ttf", g.fontDesc.FontsDir, *ufont.BoldItalic))
			boldItalicFontData, err := ioutil.ReadFile(boldItalicFontPath)
			if err != nil {
				return err
			}
			boldItalicFontHolder := NewFontHolder(boldItalicFontData)
			boldItalicFontInfo := boldItalicFontHolder.getFontInfo()
			font.BoldItalic = ufont.BoldItalic
			g.createFontAssets(boldItalicFontHolder, font, g.fontCatalog, boldItalicFontInfo.CharacterSet, boldItalicFontPath, true, true, outputPath)
		}

		g.fontCatalog.Fonts = append(g.fontCatalog.Fonts, *font)
	}
	g.createReplacementAssets(g.fontCatalog, outputPath)

	fcpath := path.Join(outputPath, fmt.Sprintf("%s_FontCatalog.json", g.fontCatalog.Name))

	data, err := g.fontCatalog.ToJson()
	if err != nil {
		return err
	}
	err = os.WriteFile(fcpath, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (g *FontCatalogGenerater) createBlockAssets(holder *FontHolder, font *Font, fontObject *FontCatalog, characterSet []rune, fontPath string, unicodeBlock *UnicodeRanges, bold bool, italic bool, outputPath string) {
	var assetSuffix string
	if bold {
		if italic {
			assetSuffix = "_BoldItalicAssets/"
		} else {
			assetSuffix = "_BoldAssets/"
		}
	} else {
		if italic {
			assetSuffix = "_ItalicAssets/"
		} else {
			assetSuffix = "_Assets/"
		}
	}
	assetsDir := path.Join(outputPath, fmt.Sprintf("%s%s", fontObject.Name, assetSuffix))
	sdfOptions := *g.opts

	sdfOptions.Filename = strings.ReplaceAll(unicodeBlock.Category, " ", "_")

	supportedCharset := ""
	for _, codePoint := range characterSet {
		if int(codePoint) >= unicodeBlock.Range[0] && int(codePoint) <= unicodeBlock.Range[1] {
			supportedCharset += string(codePoint)
		}
	}
	font.Charset += supportedCharset
	sdfOptions.Charset = supportedCharset

	if sdfOptions.Charset == "" {
		return
	} else {
		gen := NewBitmapFontGenerater(holder, sdfOptions)

		bmfont := gen.Generate()

		assetsFontDir := path.Join(assetsDir, font.Name)

		os.MkdirAll(assetsFontDir, os.ModePerm)

		for p, image := range bmfont.pageSheets {
			imagePath := path.Join(assetsDir, font.Name, fmt.Sprintf("%s.png", bmfont.Pages[p]))
			imaging.Save(image, imagePath)
		}

		font.Metrics.LineHeight = bmfont.Common.LineHeight
		font.Metrics.Base = bmfont.Common.Base

		for _, char := range bmfont.Chars {
			fontObject.MaxWidth = math.Max(fontCatalog.MaxWidth, float64(char.Width))
			fontObject.MaxHeight = math.Max(fontCatalog.MaxHeight, float64(char.Height))
		}

		data, _ := bmfont.ToJson()

		jsonPath := path.Join(assetsDir, font.Name, fmt.Sprintf("%s.json", sdfOptions.Filename))
		os.WriteFile(jsonPath, []byte(data), os.ModePerm)
	}
}

func (g *FontCatalogGenerater) createFontAssets(holder *FontHolder, font *Font, fontObject *FontCatalog, characterSet []rune, fontPath string, bold bool, italic bool, outputPath string) {
	var fontUnicodeBlockNames []string
	if len(font.Blocks) > 0 {
		fontUnicodeBlockNames = font.Blocks
	} else {
		fontUnicodeBlockNames = unicodeBlockNames
	}

	for _, blockName := range fontUnicodeBlockNames {
		var selectedBlock *UnicodeRanges
		for i := range unicodeBlocks {
			if unicodeBlocks[i].Category == blockName {
				selectedBlock = &unicodeBlocks[i]
			}
		}
		if selectedBlock == nil {
			continue
		}
		g.createBlockAssets(holder, font, fontObject, characterSet, fontPath, selectedBlock, bold, italic, outputPath)

		var blockEntry *UnicodeBlock

		for _, sb := range fontObject.SupportedBlocks {
			if sb.Name == blockName {
				blockEntry = &sb
			}
		}
		if blockEntry == nil {
			fontObject.SupportedBlocks = append(fontObject.SupportedBlocks, UnicodeBlock{
				Name:  blockName,
				Min:   selectedBlock.Range[0],
				Max:   selectedBlock.Range[1],
				Fonts: []string{font.Name},
			})
		} else if !bold && !italic {
			blockEntry.Fonts = append(blockEntry.Fonts, font.Name)
		}
	}
}

func (g *FontCatalogGenerater) createReplacementAssets(fontObject *FontCatalog, outputPath string) error {
	h := NewFontHolder([]byte(notosans_regular))
	fontInfo := h.getFontInfo()
	sdfOptions := *g.opts
	font := &Font{
		Name: "Extra",
		Metrics: FontMetric{
			Size:          sdfOptions.FontSize,
			DistanceRange: sdfOptions.DistanceRange,
			Base:          0.0,
			LineHeight:    0.0,
			LineGap:       int(math.Round(float64(fontInfo.LineGap/fontInfo.UnitsPerEm) * float64(sdfOptions.FontSize))),
			CapHeight: int(math.Round(
				float64(fontInfo.Ascent/fontInfo.UnitsPerEm) * float64(sdfOptions.FontSize))),
			XHeight: 0,
		},
		Charset: "",
	}
	assetsDir := path.Join(outputPath, fmt.Sprintf("%s%s", fontObject.Name, "_Assets/"))

	sdfOptions.Filename = "Specials"

	supportedCharset := "�"
	font.Charset += supportedCharset
	sdfOptions.Charset = supportedCharset

	gen := NewBitmapFontGenerater(h, sdfOptions)

	bmfont := gen.Generate()

	assetsFontDir := path.Join(assetsDir, "Extra")

	os.MkdirAll(assetsFontDir, os.ModePerm)

	for p, image := range bmfont.pageSheets {
		imagePath := path.Join(assetsDir, "Extra", fmt.Sprintf("%s.png", bmfont.Pages[p]))
		imaging.Save(image, imagePath)
	}

	font.Metrics.LineHeight = bmfont.Common.LineHeight
	font.Metrics.Base = bmfont.Common.Base

	for _, char := range bmfont.Chars {
		fontObject.MaxWidth = math.Max(fontCatalog.MaxWidth, float64(char.Width))
		fontObject.MaxHeight = math.Max(fontCatalog.MaxHeight, float64(char.Height))
	}

	data, _ := bmfont.ToJson()

	jsonPath := path.Join(assetsDir, "/Extra/Specials.json")
	os.WriteFile(jsonPath, []byte(data), os.ModePerm)

	var blockEntry *UnicodeBlock

	for _, sb := range fontObject.SupportedBlocks {
		if sb.Name == "Specials" {
			blockEntry = &sb
		}
	}
	if blockEntry == nil {
		fontObject.SupportedBlocks = append(fontObject.SupportedBlocks, UnicodeBlock{
			Name:  "Specials",
			Min:   65520,
			Max:   65535,
			Fonts: []string{"Extra"},
		})
	} else {
		blockEntry.Fonts = append(blockEntry.Fonts, "Extra")
	}

	g.fontCatalog.Fonts = append(g.fontCatalog.Fonts, *font)

	return nil
}
