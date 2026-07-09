package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"

	"github.com/gogpu/gg/text"
	"github.com/gogpu/gg/text/emoji"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	mdtext "github.com/yuin/goldmark/text"
)

// renderMarkdown 实现 markdown.to.image 能力：解析 markdown（goldmark，CommonMark 子集），
// 按结构渲染成 PNG。支持标题分级、加粗、斜体、行内代码、有序/无序列表、引用、代码块、分隔线、链接。
// 参数与 text.to.image 一致：context（markdown 文本）、font_color、bg_color。返回 ("image", png 字节)。
func renderMarkdown(args map[string]string) (string, []byte, error) {
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

	fonts, err := newFontSet()
	if err != nil {
		return "", nil, err
	}
	defer fonts.Close()

	src := []byte(content)
	doc := goldmark.New().Parser().Parse(mdtext.NewReader(src))

	r := &mdRenderer{
		fonts:        fonts,
		fontColor:    fontColor,
		bgColor:      backgroundColor,
		maxTextWidth: float64(defaultMaxTextWidth),
		padding:      float64(defaultPadding),
		baseSize:     float64(defaultFontSize),
	}
	r.layout(doc, src, 0)

	canvas := r.draw()

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return "", nil, fmt.Errorf("encode png: %w", err)
	}
	return "image", buf.Bytes(), nil
}

// ---- 字体集：按样式 + 字号取 face，emoji 单独 ----

type textStyle int

const (
	styleRegular textStyle = iota
	styleBold
	styleItalic
	styleCode // 行内代码：用 regular 字体（Maple Mono 本身等宽）+ 不同颜色/背景
)

type faceKey struct {
	style textStyle
	size  float64
}

// fontSet 持有三种字重的 FontSource 与 emoji 渲染器，按 (style,size) 缓存 face。
type fontSet struct {
	regularSrc *text.FontSource
	boldSrc    *text.FontSource
	italicSrc  *text.FontSource
	emoji      *colorEmojiRenderer
	faces      map[faceKey]text.Face
}

func newFontSet() (*fontSet, error) {
	regSrc, err := text.NewFontSource(mapleFontData)
	if err != nil {
		return nil, fmt.Errorf("load regular font: %w", err)
	}
	boldSrc, err := text.NewFontSource(boldFontData)
	if err != nil {
		_ = regSrc.Close()
		return nil, fmt.Errorf("load bold font: %w", err)
	}
	italicSrc, err := text.NewFontSource(italicFontData)
	if err != nil {
		_ = regSrc.Close()
		_ = boldSrc.Close()
		return nil, fmt.Errorf("load italic font: %w", err)
	}
	em, err := newColorEmojiRenderer(emojiFontData, defaultFontSize)
	if err != nil {
		_ = regSrc.Close()
		_ = boldSrc.Close()
		_ = italicSrc.Close()
		return nil, fmt.Errorf("load emoji renderer: %w", err)
	}
	return &fontSet{
		regularSrc: regSrc,
		boldSrc:    boldSrc,
		italicSrc:  italicSrc,
		emoji:      em,
		faces:      map[faceKey]text.Face{},
	}, nil
}

func (f *fontSet) Close() {
	_ = f.regularSrc.Close()
	_ = f.boldSrc.Close()
	_ = f.italicSrc.Close()
	f.emoji.Close()
}

func (f *fontSet) srcFor(style textStyle) *text.FontSource {
	switch style {
	case styleBold:
		return f.boldSrc
	case styleItalic:
		return f.italicSrc
	default: // regular / code 均用 regular（等宽）
		return f.regularSrc
	}
}

func (f *fontSet) face(style textStyle, size float64) text.Face {
	key := faceKey{style, size}
	if fc, ok := f.faces[key]; ok {
		return fc
	}
	fc := f.srcFor(style).Face(size)
	f.faces[key] = fc
	return fc
}

// advance 按样式 + 字号算单 rune 推进宽度；emoji 统一用 emoji face
func (f *fontSet) advance(rr rune, style textStyle, size float64) float64 {
	if emoji.IsEmoji(rr) {
		return f.emoji.face.Advance(string(rr))
	}
	return f.face(style, size).Advance(string(rr))
}

// measureText 测一段文本在给定样式/字号下的总宽
func (f *fontSet) measureText(s string, style textStyle, size float64) float64 {
	w := 0.0
	for _, rr := range s {
		w += f.advance(rr, style, size)
	}
	return w
}

// ---- 布局产物 ----

// textSpan inline 片段：文本 + 样式
type textSpan struct {
	text  string
	style textStyle
}

// laidLine 一条已换行的绘制单元（含块级装饰信息）
type laidLine struct {
	spans      []textSpan
	size       float64 // 该行字号（标题/代码可能不同）
	indent     float64 // 左缩进（列表/引用累积）
	bg         color.Color
	bgSet      bool
	barColor   color.Color // 引用左边框
	barSet     bool
	isHR       bool // 分隔线
	spaceAfter float64 // 段后间距
}

// ---- 渲染器 ----

const (
	mdHeadingSpace = 14
	mdParaSpace    = 12
	mdListIndent   = 30
	mdQuoteIndent  = 18
	mdCodeIndent   = 14
	mdQuoteBarW    = 4
	mdQuoteTextGap = 8
)

var (
	mdCodeColor = color.RGBA{R: 176, G: 66, B: 50, A: 255}   // 行内/块代码字色
	mdCodeBg    = color.RGBA{R: 245, G: 245, B: 245, A: 255} // 代码背景浅灰
	mdQuoteBar  = color.RGBA{R: 180, G: 180, B: 180, A: 255} // 引用左边框灰
)

type mdRenderer struct {
	fonts        *fontSet
	fontColor    color.Color
	bgColor      color.Color
	maxTextWidth float64
	padding      float64
	baseSize     float64

	lines []laidLine
}

// layout 遍历 block 级节点，生成 laidLine 追加到 r.lines。indent 为累积左缩进。
func (r *mdRenderer) layout(n gast.Node, src []byte, indent float64) {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		r.layoutBlock(child, src, indent)
	}
}

func (r *mdRenderer) layoutBlock(n gast.Node, src []byte, indent float64) {
	switch n.Kind() {
	case gast.KindHeading:
		r.layoutHeading(n, src, indent)
	case gast.KindParagraph, gast.KindTextBlock:
		// tight list（项间无空行）的内容是 TextBlock，loose list / 普通段落是 Paragraph；
		// 两者 inline 处理一致，共用 layoutParagraph
		r.layoutParagraph(n, src, indent)
	case gast.KindList:
		r.layoutList(n, src, indent)
	case gast.KindBlockquote:
		r.layoutBlockquote(n, src, indent)
	case gast.KindFencedCodeBlock, gast.KindCodeBlock:
		r.layoutCode(n, src, indent)
	case gast.KindThematicBreak:
		r.lines = append(r.lines, laidLine{size: r.baseSize, indent: indent, isHR: true, spaceAfter: mdParaSpace})
	case gast.KindHTMLBlock:
		// HTML 块按纯文本段落兜底（用 Lines() 取原始文本；HTMLBlock.Text 已弃用）
		var buf strings.Builder
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			buf.Write(seg.Value(src))
		}
		spans := []textSpan{{text: buf.String(), style: styleRegular}}
		r.appendWrapped(spans, r.baseSize, indent, mdParaSpace)
	}
}

func (r *mdRenderer) layoutHeading(n gast.Node, src []byte, indent float64) {
	level := 1
	if h, ok := n.(*gast.Heading); ok {
		level = h.Level
	}
	size := headingSize(level)
	spans := collectSpans(n, src, styleBold)
	// 标题整体加粗：强制 bold（保留内部 code 样式）
	for i := range spans {
		if spans[i].style != styleCode {
			spans[i].style = styleBold
		}
	}
	r.appendWrapped(spans, size, indent, mdHeadingSpace)
}

func (r *mdRenderer) layoutParagraph(n gast.Node, src []byte, indent float64) {
	spans := collectSpans(n, src, styleRegular)
	r.appendWrapped(spans, r.baseSize, indent, mdParaSpace)
}

func (r *mdRenderer) layoutList(n gast.Node, src []byte, indent float64) {
	list, _ := n.(*gast.List)
	idx := 1
	for item := n.FirstChild(); item != nil; item = item.NextSibling() {
		marker := "•  "
		if list != nil && list.IsOrdered() {
			marker = fmt.Sprintf("%d. ", idx)
			idx++
		}
		itemStart := len(r.lines)
		for child := item.FirstChild(); child != nil; child = child.NextSibling() {
			// 列表项内嵌套块（含嵌套列表）按缩进递增
			r.layoutBlock(child, src, indent+mdListIndent)
		}
		// 首行前缀项目符号
		if itemStart < len(r.lines) {
			first := &r.lines[itemStart]
			first.spans = append([]textSpan{{text: marker, style: styleRegular}}, first.spans...)
		}
	}
}

func (r *mdRenderer) layoutBlockquote(n gast.Node, src []byte, indent float64) {
	start := len(r.lines)
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		r.layoutBlock(child, src, indent+mdQuoteIndent)
	}
	for i := start; i < len(r.lines); i++ {
		r.lines[i].barSet = true
		r.lines[i].barColor = mdQuoteBar
	}
}

func (r *mdRenderer) layoutCode(n gast.Node, src []byte, indent float64) {
	var codeLines []string
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		codeLines = append(codeLines, strings.TrimRight(string(seg.Value(src)), "\n"))
	}
	if len(codeLines) == 0 {
		return
	}
	for _, cl := range codeLines {
		r.lines = append(r.lines, laidLine{
			spans:  []textSpan{{text: cl, style: styleCode}},
			size:   r.baseSize,
			indent: indent + mdCodeIndent,
			bg:     mdCodeBg, bgSet: true,
		})
	}
	r.lines[len(r.lines)-1].spaceAfter = mdParaSpace
}

// appendWrapped 把 inline 片段按宽度换行后追加，末行设段后间距
func (r *mdRenderer) appendWrapped(spans []textSpan, size, indent, spaceAfter float64) {
	wrapped := r.wrapSpans(spans, size, indent)
	for i := range wrapped {
		wrapped[i].size = size
	}
	if len(wrapped) > 0 {
		wrapped[len(wrapped)-1].spaceAfter = spaceAfter
	}
	r.lines = append(r.lines, wrapped...)
}

// wrapSpans 按 rune 与各 span 样式对应的 face 宽度换行。indent 计入可用宽度扣减。
func (r *mdRenderer) wrapSpans(spans []textSpan, size, indent float64) []laidLine {
	maxWidth := r.maxTextWidth - indent
	if maxWidth < r.baseSize {
		maxWidth = r.baseSize
	}
	var lines []laidLine
	var cur []textSpan
	curW := 0.0

	flush := func() {
		lines = append(lines, laidLine{spans: cur, indent: indent})
		cur = nil
		curW = 0
	}

	pushRune := func(rr rune, style textStyle) {
		if n := len(cur); n > 0 && cur[n-1].style == style {
			cur[n-1].text += string(rr)
		} else {
			cur = append(cur, textSpan{text: string(rr), style: style})
		}
		curW += r.fonts.advance(rr, style, size)
	}

	for _, sp := range spans {
		for _, rr := range sp.text {
			w := r.fonts.advance(rr, sp.style, size)
			if curW+w > maxWidth && len(cur) > 0 {
				flush()
				if rr == ' ' {
					continue // 行首省略空格
				}
			}
			pushRune(rr, sp.style)
		}
	}
	if len(cur) > 0 || len(lines) == 0 {
		flush()
	}
	return lines
}

// collectSpans 把 inline 容器（段落/标题的内容）递归展平为 []textSpan，保留样式继承。
func collectSpans(n gast.Node, src []byte, inherit textStyle) []textSpan {
	var spans []textSpan
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case gast.KindText:
			t, ok := child.(*gast.Text)
			if !ok {
				continue
			}
			spans = append(spans, textSpan{text: string(t.Value(src)), style: inherit})
			// goldmark 把段内软/硬换行作为 Text 节点的标志位（非独立节点）：
			// 追加空格避免相邻文本粘连
			if t.SoftLineBreak() || t.HardLineBreak() {
				spans = append(spans, textSpan{text: " ", style: inherit})
			}
		case gast.KindCodeSpan:
			spans = append(spans, collectSpans(child, src, styleCode)...)
		case gast.KindEmphasis:
			// goldmark 用 Emphasis + Level 区分：Level=2 为 **加粗**，Level=1 为 *斜体*
			style := styleItalic
			if em, ok := child.(*gast.Emphasis); ok && em.Level >= 2 {
				style = styleBold
			}
			spans = append(spans, collectSpans(child, src, style)...)
		case gast.KindLink:
			spans = append(spans, collectSpans(child, src, inherit)...) // 仅保留链接文字
		}
	}
	return spans
}

// lineHeight 单行高度：主字体与 emoji 取大者
func (r *mdRenderer) lineHeight(line laidLine) float64 {
	f := r.fonts.face(styleRegular, line.size)
	return math.Max(f.Metrics().LineHeight(), r.fonts.emoji.face.Metrics().LineHeight()) * defaultLineSpacing
}

func (r *mdRenderer) measureLine(line laidLine) float64 {
	if line.isHR {
		return r.maxTextWidth
	}
	w := 0.0
	for _, sp := range line.spans {
		w += r.fonts.measureText(sp.text, sp.style, line.size)
	}
	return w
}

// draw 测总算画布尺寸，建画布后逐行绘制
func (r *mdRenderer) draw() *image.RGBA {
	var totalHeight, maxW float64
	for _, line := range r.lines {
		totalHeight += r.lineHeight(line) + line.spaceAfter
		if lw := line.indent + r.measureLine(line); lw > maxW {
			maxW = lw
		}
	}

	canvasW := int(math.Ceil(r.padding*2 + maxW))
	canvasH := int(math.Ceil(r.padding*2 + totalHeight))
	if canvasW < 1 {
		canvasW = 1
	}
	if canvasH < 1 {
		canvasH = 1
	}
	canvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(r.bgColor), image.Point{}, draw.Src)

	y := r.padding
	for _, line := range r.lines {
		lh := r.lineHeight(line)
		lineX := r.padding + line.indent

		if line.bgSet {
			bgW := float64(canvasW) - r.padding - lineX
			fillRect(canvas, lineX, y, bgW, lh, line.bg)
		}
		if line.isHR {
			fillRect(canvas, lineX, y+lh/2-1, r.maxTextWidth, 2, r.fontColor)
			y += lh + line.spaceAfter
			continue
		}
		if line.barSet {
			fillRect(canvas, lineX, y, mdQuoteBarW, lh, line.barColor)
		}
		textX := lineX
		if line.barSet {
			textX += mdQuoteTextGap
		}
		baseline := y + r.fonts.face(styleRegular, line.size).Metrics().Ascent
		r.drawSpans(canvas, line.spans, line.size, textX, baseline)

		y += lh + line.spaceAfter
	}
	return canvas
}

// drawSpans 按各 span 样式选 face/颜色绘制，支持段内 face 切换；emoji 段用 emoji renderer
func (r *mdRenderer) drawSpans(canvas draw.Image, spans []textSpan, size, x, baseline float64) {
	for _, sp := range spans {
		face := r.fonts.face(sp.style, size)
		col := r.colorFor(sp.style)
		for _, run := range emoji.Segment(sp.text) {
			if run.IsEmoji {
				for _, rr := range run.Text {
					r.fonts.emoji.DrawRune(canvas, rr, x, baseline)
					x += r.fonts.emoji.face.Advance(string(rr))
				}
				continue
			}
			text.Draw(canvas, run.Text, face, x, baseline, col)
			x += face.Advance(run.Text)
		}
	}
}

func (r *mdRenderer) colorFor(style textStyle) color.Color {
	if style == styleCode {
		return mdCodeColor
	}
	return r.fontColor
}

// ---- 辅助 ----

func headingSize(level int) float64 {
	switch level {
	case 1:
		return 38
	case 2:
		return 32
	case 3:
		return 28
	case 4:
		return 26
	default:
		return 24
	}
}

// fillRect 在画布上填充矩形（用于代码背景、引用边框、分隔线）
func fillRect(dst draw.Image, x, y, w, h float64, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	rect := image.Rect(int(x), int(y), int(x+w), int(y+h))
	draw.Draw(dst, rect, image.NewUniform(c), image.Point{}, draw.Src)
}
