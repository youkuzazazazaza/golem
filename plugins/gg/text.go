package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"strings"

	"github.com/gogpu/gg/text"
)

// text.to.image 布局参数
const (
	defaultMaxTextWidth = 1600
	defaultLineSpacing  = 0.8
	defaultPadding      = 25
)

// renderText 实现 text.to.image 能力：把纯文本按字体宽度自动换行渲染成 PNG。
// 参数：context（必填，纯文本）、font_color、bg_color（可选）。返回 ("image", png 字节)。
// 注意：本能力不解析 markdown，传入的文本原样绘制；markdown 渲染见 markdown.to.image。
func renderText(args map[string]string) (string, []byte, error) {
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
