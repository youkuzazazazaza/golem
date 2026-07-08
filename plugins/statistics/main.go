package main

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/sbgayhub/golem/sdk/chatroom"
	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)

func main() {
	plugin.Start(&StatisticsPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: defaultConfig(),
		},
	})
}

type StatisticsPlugin struct {
	plugin.ConfigAbility[Config]

	message  message.Ability
	chatroom chatroom.Ability
	contact  contact.Ability
	caller   plugin.CallerAbility

	mu    sync.Mutex
	dbDir string
	store *store
}

func (p *StatisticsPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "statistics",
		Author:      "ovo",
		Version:     "1.1.0",
		Description: "消息统计 + 群成员人物画像插件：记录消息、提供发言排行/详情，并支持「人物画像 <昵称>」基于历史发言生成画像（复用 ai 插件 ai.chat 能力）",
		Priority:    -1 << 31,
		Next:        false,
		AlwaysRun:   true,
	}
}

func (p *StatisticsPlugin) GetSubscriptions() []string {
	return []string{"message"}
}

func (p *StatisticsPlugin) OnLoad() error {
	return p.ensureStore()
}

func (p *StatisticsPlugin) OnUnload() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.store == nil {
		return nil
	}
	err := p.store.Close()
	p.store = nil
	return err
}

func (p *StatisticsPlugin) OnEnable() error {
	return p.ensureStore()
}

func (p *StatisticsPlugin) OnDisable() error {
	return nil
}

func (p *StatisticsPlugin) OnEvent(event *plugin.Event) (bool, error) {
	msg := event.GetPayload().(*plugin.Event_Message).Message
	if msg.GetSender() == nil {
		return false, nil
	}

	// 排行关键词：群聊精确匹配，不记录，直接处理（加锁，DB 查询快）
	if isRankingKeyword(msg) {
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.store == nil {
			return false, errors.New("store is not initialized")
		}
		return p.handleRanking(msg)
	}

	// 先记录消息（加锁，写入快）
	p.mu.Lock()
	if p.store == nil {
		p.mu.Unlock()
		return false, errors.New("store is not initialized")
	}
	_, err := p.store.record(msg)
	p.mu.Unlock()
	if err != nil {
		slog.Warn("[statistics] 记录消息失败", "err", err)
		return false, nil
	}

	// 人物画像触发：生成（不加锁，AI 调用慢）+ 消费事件，避免 ai 重复回复
	if name, global, rebuild, triggered := parseTrigger(msg); triggered {
		return p.handleProfile(msg, name, global, rebuild)
	}

	// 普通消息：已记录，放行给后续插件（如 ai）。Next=false 但 handled=false 仍会继续分发。
	return false, nil
}

func (p *StatisticsPlugin) ensureStore() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.store != nil {
		return nil
	}

	st, err := openStore(p.dbDir)
	if err != nil {
		return err
	}
	p.store = st
	return nil
}

func (p *StatisticsPlugin) sendText(receiver *contact.Contact, content string, reminds []string) error {
	if p.message == nil {
		return errors.New("message ability is not injected")
	}
	if receiver == nil || receiver.GetUsername() == "" {
		return errors.New("message receiver is empty")
	}

	_, err := p.message.Send(&message.Message{
		Type:     message.TypeText,
		Receiver: receiver,
		Content:  content,
		Data: &message.Message_Text{Text: &message.TextData{
			Content: content,
			Reminds: reminds,
		}},
	})
	return err
}

// sendProfileResult 按配置决定把画像渲染成图片（经 gg 插件 text.to.image）还是发纯文本。
// render_image=true 时尝试渲染图片；gg 未启用/渲染失败自动回退文本。render_image=false 直接发文本。
func (p *StatisticsPlugin) sendProfileResult(receiver *contact.Contact, text string) error {
	if normalizeConfig(p.Config).RenderImage {
		if data, err := p.renderTextToImage(text); err == nil && len(data) > 0 {
			if imgErr := p.sendImage(receiver, data); imgErr == nil {
				return nil
			}
		} else {
			slog.Warn("[statistics] 画像渲染图片失败，回退为文本", "err", err)
		}
	}
	return p.sendText(receiver, text, nil)
}

// renderTextToImage 调用 text.to.image 能力把文本渲���成 PNG
func (p *StatisticsPlugin) renderTextToImage(text string) ([]byte, error) {
	if p.caller == nil {
		return nil, fmt.Errorf("caller 未注入（需要 gg 插件提供 text.to.image）")
	}
	mime, data, err := p.caller.CallPlugin("text.to.image", map[string]string{"context": text})
	if err != nil {
		return nil, err
	}
	_ = mime
	if len(data) == 0 {
		return nil, fmt.Errorf("text.to.image 返回空数据")
	}
	return data, nil
}

// sendImage 发送图片消息（PNG 字节）
func (p *StatisticsPlugin) sendImage(receiver *contact.Contact, pngData []byte) error {
	if p.message == nil {
		return errors.New("message ability is not injected")
	}
	if receiver == nil || receiver.GetUsername() == "" {
		return errors.New("message receiver is empty")
	}
	_, err := p.message.Send(&message.Message{
		Type:     message.TypeImage,
		Receiver: receiver,
		Content:  "人物画像",
		Data: &message.Message_Image{Image: &message.ImageData{
			Media: &message.Media{Data: pngData},
		}},
	})
	return err
}
