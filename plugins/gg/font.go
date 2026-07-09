package main

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"

	xdraw "golang.org/x/image/draw"

	"github.com/gogpu/gg/text"
	"github.com/gogpu/gg/text/emoji"
)

// 字体资源（Regular/Bold/Italic 供主字体使用，Emoji 供彩色 emoji 渲染）。
// Bold/Italic 为 markdown.to.image 的加粗/斜体而加载；text.to.image 仅用 Regular。
//
//go:embed MapleMono-NF-CN-Regular.ttf
var mapleFontData []byte

//
//go:embed MapleMono-NF-CN-Bold.ttf
var boldFontData []byte

//
//go:embed MapleMono-NF-CN-Italic.ttf
var italicFontData []byte

//
//go:embed NotoColorEmoji.ttf
var emojiFontData []byte

const defaultFontSize = 24

// colorEmojiRenderer emoji 感知的字体渲染基础设施。
// 主字体（Regular/Bold/Italic）以 Face 形式由调用方传入；
// 本结构负责 emoji 的 CBDT/CBLC 位图提取与绘制，以及基于主字体的文本测量/换行/绘制。
// text.to.image 与 markdown.to.image 共用此能力。
type colorEmojiRenderer struct {
	source    *text.FontSource
	face      text.Face
	extractor *emoji.CBDTExtractor
	ppem      uint16
	size      float64
}

func newColorEmojiRenderer(fontData []byte, size float64) (*colorEmojiRenderer, error) {
	source, err := text.NewFontSource(fontData)
	if err != nil {
		return nil, err
	}

	extractor, err := emoji.NewCBDTExtractor(getTable(fontData, "CBDT"), getTable(fontData, "CBLC"))
	if err != nil {
		_ = source.Close()
		return nil, err
	}

	ppems := extractor.AvailablePPEMs()
	if len(ppems) == 0 {
		_ = source.Close()
		return nil, fmt.Errorf("emoji font has no bitmap strikes")
	}

	return &colorEmojiRenderer{
		source:    source,
		face:      source.Face(size),
		extractor: extractor,
		ppem:      ppems[len(ppems)-1],
		size:      size,
	}, nil
}

func (r *colorEmojiRenderer) Close() {
	_ = r.source.Close()
}

func (r *colorEmojiRenderer) WrapText(mainFace text.Face, s string, maxWidth float64) []string {
	var lines []string
	for _, paragraph := range strings.Split(s, "\n") {
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}

		lines = append(lines, r.wrapParagraph(mainFace, paragraph, maxWidth)...)
	}
	return lines
}

func (r *colorEmojiRenderer) wrapParagraph(mainFace text.Face, s string, maxWidth float64) []string {
	var lines []string
	var current strings.Builder
	currentWidth := 0.0

	for _, rr := range s {
		charWidth := r.runeAdvance(mainFace, rr)
		if current.Len() > 0 && currentWidth+charWidth > maxWidth {
			lines = append(lines, current.String())
			current.Reset()
			currentWidth = 0
		}

		if current.Len() == 0 && rr == ' ' {
			continue
		}

		current.WriteRune(rr)
		currentWidth += charWidth
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return lines
}

func (r *colorEmojiRenderer) DrawLines(dst draw.Image, mainFace text.Face, lines []string, x, y float64, lineSpacing float64, col color.Color) {
	lineHeight := r.LineHeight(mainFace, lineSpacing)
	baseline := y + mainFace.Metrics().Ascent

	for _, line := range lines {
		if line != "" {
			r.DrawLine(dst, mainFace, line, x, baseline, col)
		}
		baseline += lineHeight
	}
}

func (r *colorEmojiRenderer) MeasureLines(mainFace text.Face, lines []string) float64 {
	maxWidth := 0.0
	for _, line := range lines {
		width := r.MeasureLine(mainFace, line)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (r *colorEmojiRenderer) MeasureLine(mainFace text.Face, s string) float64 {
	width := 0.0
	for _, rr := range s {
		width += r.runeAdvance(mainFace, rr)
	}
	return width
}

func (r *colorEmojiRenderer) DrawLine(dst draw.Image, mainFace text.Face, s string, x, y float64, col color.Color) {
	for _, run := range emoji.Segment(s) {
		if run.IsEmoji {
			for _, rr := range run.Text {
				r.DrawRune(dst, rr, x, y)
				x += r.face.Advance(string(rr))
			}
			continue
		}

		text.Draw(dst, run.Text, mainFace, x, y, col)
		x += mainFace.Advance(run.Text)
	}
}

func (r *colorEmojiRenderer) LineHeight(mainFace text.Face, lineSpacing float64) float64 {
	return math.Max(mainFace.Metrics().LineHeight(), r.face.Metrics().LineHeight()) * lineSpacing
}

func (r *colorEmojiRenderer) runeAdvance(mainFace text.Face, rr rune) float64 {
	if emoji.IsEmoji(rr) {
		return r.face.Advance(string(rr))
	}
	return mainFace.Advance(string(rr))
}

func (r *colorEmojiRenderer) DrawRune(dst draw.Image, rr rune, x, y float64) {
	gid := r.source.Parsed().GlyphIndex(rr)
	glyph, err := r.extractor.GetGlyph(gid, r.ppem)
	if err != nil {
		return
	}

	img, err := png.Decode(bytes.NewReader(glyph.Data))
	if err != nil {
		return
	}

	scale := r.size / float64(glyph.PPEM)
	if scale == 0 {
		scale = 1
	}

	dstX := int(x + float64(glyph.OriginX)*scale)
	dstY := int(y - float64(glyph.OriginY)*scale)
	dstW := int(float64(glyph.Width) * scale)
	dstH := int(float64(glyph.Height) * scale)
	if dstW <= 0 || dstH <= 0 {
		return
	}

	rect := image.Rect(dstX, dstY, dstX+dstW, dstY+dstH)
	xdraw.CatmullRom.Scale(dst, rect, img, img.Bounds(), xdraw.Over, nil)
}

// getTable 从原始字体字节中按 tag（如 "CBDT"/"CBLC"）取出对应表数据
func getTable(data []byte, tag string) []byte {
	if len(data) < 12 {
		return nil
	}

	numTables := int(binary.BigEndian.Uint16(data[4:6]))
	offset := 12
	for i := 0; i < numTables && offset+16 <= len(data); i++ {
		t := string(data[offset : offset+4])
		tableOffset := binary.BigEndian.Uint32(data[offset+8 : offset+12])
		tableLength := binary.BigEndian.Uint32(data[offset+12 : offset+16])
		if t == tag && tableOffset+tableLength <= uint32(len(data)) {
			return data[tableOffset : tableOffset+tableLength]
		}
		offset += 16
	}
	return nil
}
