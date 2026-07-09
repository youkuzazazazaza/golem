package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// parseBackgroundColor 解析背景色参数；空字符串返回默认浅灰 #F7F7F7（与 host 内置命令图一致）。
// 必须为不透明色：透明背景会让 Go png encoder 输出带 alpha 的 RGBA（color type 6），
// 微信缩略图生成器会将其当作黑底预乘，导致缩略图发黑（大图正常）。
func parseBackgroundColor(s string) (color.Color, error) {
	if strings.TrimSpace(s) == "" {
		return color.RGBA{R: 0xF7, G: 0xF7, B: 0xF7, A: 0xFF}, nil
	}
	return parseFontColor(s)
}

// parseFontColor 解析前景色参数；支持颜色名（black/white/red/...）与 #HEX / 0xHEX（3/4/6/8 位）。
// 空字符串返回黑色（默认前景色）。
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

// expandShortHex 把 3/4 位短 HEX 展开为 8 位（含 alpha）
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
