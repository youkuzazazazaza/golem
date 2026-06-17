package main

import (
	"log/slog"
	"math/rand/v2"
	"sync"

	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)

// AiPlugin AI 插件主结构
type AiPlugin struct {
	plugin.ConfigAbility[Config]
	contact contact.Ability
	message message.Ability

	configMu        sync.RWMutex
	sessionConfigMu sync.RWMutex
	selfMu          sync.RWMutex
	self            *contact.SelfInfo
	owner           *contact.Contact
	mu              sync.Mutex
	sessions        map[string][]openAIMessage
}

// Config 插件配置
type Config struct {
	BaseURL            string                    `toml:"base_url" comment:"OpenAI 兼容接口地址，例如 https://api.openai.com/v1"`
	APIKey             string                    `toml:"api_key" comment:"接口密钥"`
	Model              string                    `toml:"model" comment:"模型名称"`
	ActivePrompt       string                    `toml:"active_prompt" comment:"当前使用的提示词名称"`
	Prompts            map[string]string         `toml:"prompts" comment:"提示词映射，key 为提示词名称"`
	LegacyPrompt       string                    `toml:"prompt,omitempty" comment:"旧版提示词配置，启动后迁移到 prompts.default"`
	ReplyRate          float64                   `toml:"reply_rate" comment:"普通消息回复概率，取值 0~1"`
	MaxContextMessages int                       `toml:"max_context_messages" comment:"每个会话最多保留的上下文消息数"`
	HTTPTimeoutSeconds int                       `toml:"http_timeout_seconds" comment:"大模型请求超时时间，单位秒"`
	SessionConfigs     map[string]*SessionConfig `toml:"session_configs,omitempty" comment:"会话级配置，key 为会话标识"`
}

// newAiPlugin 创建 AI 插件实例
func newAiPlugin() (*AiPlugin, error) {
	p := &AiPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: defaultConfig(),
		},
		sessions: map[string][]openAIMessage{},
	}
	if err := registerCommands(p); err != nil {
		return nil, err
	}
	return p, nil
}

// shouldReply 判断是否应该回复
func (p *AiPlugin) shouldReply(in incomingMessage) bool {
	if in.MentionedBot || in.QuotedBot {
		return true
	}
	if !in.IsChatroom {
		return true
	}
	rate := p.getReplyRate(in.SessionKey)
	if rate <= 0 {
		return false
	}
	if rate >= 1 {
		return true
	}
	return rand.Float64() < rate
}

// main 插件入口
func main() {
	p, err := newAiPlugin()
	if err != nil {
		slog.Error("[ai] 初始化失败", "err", err)
		return
	}
	plugin.Start(p)
}
