package main

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)

type UniversalPlugin struct {
	plugin.ConfigAbility[Config]

	message message.Ability

	mu             sync.RWMutex
	rulesByKeyword map[string]*Rule
}

type Config struct {
	Rules              []Rule `toml:"rules" comment:"规则列表"`
	HTTPTimeoutSeconds int    `toml:"http_timeout_seconds" comment:"HTTP 请求超时时间，单位秒"`
}

func (p *UniversalPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "universal",
		Author:      "ovo",
		Version:     "1.0.0",
		Description: "规则驱动的通用 API 请求与消息发送插件",
		Next:        false,
		Priority:    0,
		AlwaysRun:   false,
	}
}

func (p *UniversalPlugin) GetSubscriptions() []string {
	return []string{"message::text", "message::app::quote"}
}

func (p *UniversalPlugin) GetCommands() []string {
	return plugin.CommandCommands()
}

func (p *UniversalPlugin) OnCommand(cmd *plugin.Command) (string, error) {
	return plugin.DispatchCommand(cmd)
}

func (p *UniversalPlugin) OnLoad() error {
	return p.rebuildIndex()
}

func (p *UniversalPlugin) OnUnload() error {
	return nil
}

func (p *UniversalPlugin) OnEnable() error {
	return p.rebuildIndex()
}

func (p *UniversalPlugin) OnDisable() error {
	return nil
}

func (p *UniversalPlugin) OnEvent(e *plugin.Event) (bool, error) {
	payload, ok := e.Payload.(*plugin.Event_Message)
	if !ok || payload.Message == nil {
		return false, nil
	}

	parsed, ok := parseIncomingText(messageContent(payload.Message))
	if !ok {
		return false, nil
	}

	rule, ok := p.ruleForKeyword(parsed.Keyword)
	if !ok {
		return false, nil
	}

	quote := extractQuote(payload.Message)
	vars := buildTemplateVars(parsed, quote)
	result, err := executeRule(rule, vars, p.requestTimeout())
	if err != nil {
		return true, err
	}

	receiver := e.GetSender()
	if receiver == "" && payload.Message.GetSender() != nil {
		receiver = payload.Message.GetSender().GetUsername()
	}
	if receiver == "" {
		return true, errors.New("message receiver is empty")
	}
	if err := p.sendResult(payload.Message.Sender, rule.SendType, result, collectMentionTargets(rule, parsed, quote, payload.Message)); err != nil {
		return true, err
	}
	return true, nil
}

func (p *UniversalPlugin) requestTimeout() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()

	seconds := p.Config.HTTPTimeoutSeconds
	if seconds <= 0 {
		seconds = defaultHTTPTimeoutSeconds
	}
	return time.Duration(seconds) * time.Second
}

func main() {
	p := &UniversalPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: Config{
				HTTPTimeoutSeconds: defaultHTTPTimeoutSeconds,
			},
		},
		rulesByKeyword: map[string]*Rule{},
	}
	if err := p.registerCommands(); err != nil {
		slog.Error("注册命令失败", "err", err)
		return
	}
	plugin.Start(p)
}
