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
	"strconv"
	"strings"

	xdraw "golang.org/x/image/draw"

	"github.com/gogpu/gg/text"
	"github.com/gogpu/gg/text/emoji"
	"github.com/sbgayhub/golem/sdk/plugin"
)

const (
	defaultMaxTextWidth = 1600
	defaultFontSize     = 24
	defaultLineSpacing  = 0.8
	defaultPadding      = 25
)

//go:embed MapleMono-NF-CN-Regular.ttf
var mapleFontData []byte

//go:embed NotoColorEmoji.ttf
var emojiFontData []byte

type GGPlugin struct {
}

func (g *GGPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "gg",
		Author:      "ovo",
		Version:     "v0.0.0",
		Description: "使用gogpu/gg库生成图片",
		Priority:    0,
		Next:        false,
		AlwaysRun:   false,
	}
}

func (g *GGPlugin) GetCapabilities() []string {
	return []string{"text.to.image"}
}

func (g *GGPlugin) OnCall(capability string, args map[string]string) (string, []byte, error) {
	if capability != "text.to.image" {
		return "", nil, fmt.Errorf("unsupported capability: %s", capability)
	}

	content := strings.TrimSpace(args["context"])
	if content == "" {
		return "", nil, fmt.Errorf("context is required")
	}

	fontColor, err := parseFontColor(args["font_color"])
	if err != nil {
		return "", nil, err
	}

	backgroundColor, err := parseBackgroundColor(args["bg_color"])
	if err != nil {
		return "", nil, err
	}

	mainSource, err := text.NewFontSource(mapleFontData)
	if err != nil {
		return "", nil, fmt.Errorf("load main font: %w", err)
	}
	defer func() { _ = mainSource.Close() }()

	mainFace := mainSource.Face(defaultFontSize)

	emojiRenderer, err := newColorEmojiRenderer(emojiFontData, defaultFontSize)
	if err != nil {
		return "", nil, fmt.Errorf("load emoji renderer: %w", err)
	}
	defer emojiRenderer.Close()

	maxTextWidth := float64(defaultMaxTextWidth)
	lines := emojiRenderer.WrapText(mainFace, content, maxTextWidth)
	lineHeight := emojiRenderer.LineHeight(mainFace, defaultLineSpacing)
	maxLineWidth := emojiRenderer.MeasureLines(mainFace, lines)

	canvasWidth := defaultPadding*2 + int(math.Ceil(maxLineWidth))
	canvasHeight := defaultPadding*2 + int(math.Ceil(float64(len(lines))*lineHeight))

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)

	emojiRenderer.DrawLines(canvas, mainFace, lines, defaultPadding, defaultPadding, defaultLineSpacing, fontColor)

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return "", nil, fmt.Errorf("encode png: %w", err)
	}

	return "image", buf.Bytes(), nil
}

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

func parseBackgroundColor(s string) (color.Color, error) {
	if strings.TrimSpace(s) == "" {
		return color.Transparent, nil
	}
	return parseFontColor(s)
}

func parseFontColor(s string) (color.Color, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return color.Black, nil
	}

	switch strings.ToLower(s) {
	case "black":
		return color.Black, nil
	case "white":
		return color.White, nil
	case "red":
		return color.RGBA{R: 255, A: 255}, nil
	case "green":
		return color.RGBA{G: 255, A: 255}, nil
	case "blue":
		return color.RGBA{B: 255, A: 255}, nil
	}

	s = strings.TrimPrefix(s, "#")
	s = strings.TrimPrefix(strings.ToLower(s), "0x")

	switch len(s) {
	case 3:
		s = expandShortHex(s, false)
	case 4:
		s = expandShortHex(s, true)
	case 6:
		s += "ff"
	case 8:
	default:
		return nil, fmt.Errorf("invalid font_color: %q", s)
	}

	value, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid font_color: %w", err)
	}

	return color.RGBA{
		R: uint8(value >> 24),
		G: uint8(value >> 16),
		B: uint8(value >> 8),
		A: uint8(value),
	}, nil
}

func expandShortHex(s string, withAlpha bool) string {
	var b strings.Builder
	for _, rr := range s {
		b.WriteRune(rr)
		b.WriteRune(rr)
	}
	if !withAlpha {
		b.WriteString("ff")
	}
	return b.String()
}

func main() {
	plugin.Start(&GGPlugin{})
}
