package main

import (
	"fmt"

	"github.com/sbgayhub/golem/sdk/plugin"
)

// GGPlugin 图片生成插件：提供文本/Markdown 渲染成 PNG 的能力。
// - text.to.image：纯文本按字体宽度换行渲染（实现见 text.go）
// - markdown.to.image：markdown 解析后按结构渲染（实现见 markdown.go）
type GGPlugin struct {
}

// GetMetadata 返回插件元数据
func (g *GGPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "gg",
		Author:      "ovo",
		Version:     "v0.0.0",
		Description: "使用gogpu/gg库生成图片，提供 text.to.image 与 markdown.to.image 能力",
		Priority:    0,
		Next:        false,
		AlwaysRun:   false,
	}
}

// GetCapabilities 声明本插件可被其他插件调用的能力
func (g *GGPlugin) GetCapabilities() []string {
	return []string{"text.to.image", "markdown.to.image"}
}

// OnCall 按能力名路由到对应渲染实现；新增能力在此 switch 加分支即可，
// 既有能力的实现互不干扰。
func (g *GGPlugin) OnCall(capability string, args map[string]string) (string, []byte, error) {
	switch capability {
	case "text.to.image":
		return renderText(args)
	case "markdown.to.image":
		return renderMarkdown(args)
	default:
		return "", nil, fmt.Errorf("unsupported capability: %s", capability)
	}
}

func main() {
	plugin.Start(&GGPlugin{})
}
