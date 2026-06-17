package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sbgayhub/golem/sdk/plugin"
)

type aiGetCommand struct {
	_ struct{} `cmd:"ai get" help:"查看 AI 配置" usage:"/ai get" example:"/ai get"`
}

type aiSetCommand struct {
	_          struct{} `cmd:"ai set" help:"更新 AI 配置" usage:"/ai set [-u url] [-k api_key] [-m model] [-r reply_rate] [--max-context n] [--timeout n]" example:"/ai set -u https://api.deepseek.com/v1 -k sk-xxx -m deepseek-chat\n/ai set -r 0.2"`
	BaseURL    *string  `flag:"u,url" help:"OpenAI 兼容接口地址，例如 https://api.openai.com/v1"`
	APIKey     *string  `flag:"k,api-key" help:"接口密钥"`
	Model      *string  `flag:"m,model" help:"模型名称"`
	ReplyRate  *float64 `flag:"r,reply-rate" help:"普通消息回复概率，取值 0~1"`
	MaxContext *int     `flag:"max-context" help:"每个会话最多保留的上下文消息数"`
	Timeout    *int     `flag:"timeout" help:"大模型请求超时时间，单位秒"`
}

type aiPromptCommand struct {
	_       struct{} `cmd:"ai prompt" help:"切换或设置提示词" usage:"/ai prompt <name> [-p prompt]" example:"/ai prompt default\n/ai prompt roleplay -p \"你是一个微信聊天机器人\""`
	Name    string   `arg:"name" help:"提示词名称" required:"true"`
	Prompt  *string  `flag:"p,prompt" help:"提示词内容；提供后会新增或更新该提示词并切换到它"`
	Command *plugin.Command
}

type aiClearContextCommand struct {
	_       struct{} `cmd:"ai clear-context" help:"清理上下文" usage:"/ai clear-context [-t target]" example:"/ai clear-context\n/ai clear-context -t chatroom::123@chatroom"`
	Target  string   `flag:"t,target" help:"会话 key；不传则清理当前命令来源会话"`
	Command *plugin.Command
}

type aiSessionGetCommand struct {
	_       struct{} `cmd:"ai session-get" help:"查看会话配置" usage:"/ai session-get [-t target]" example:"/ai session-get\n/ai session-get -t chatroom:123@chatroom"`
	Target  string   `flag:"t,target" help:"会话 key；不传则查看当前会话"`
	Command *plugin.Command
}

type aiSessionSetCommand struct {
	_          struct{} `cmd:"ai session-set" help:"设置会话配置" usage:"/ai session-set [-t target] [-r reply_rate] [-p prompt] [-c max_context]" example:"/ai session-set -r 0.8 -p meow\n/ai session-set -t chatroom:123@chatroom -r 1.0"`
	Target     string   `flag:"t,target" help:"会话 key；不传则设置当前会话"`
	ReplyRate  *float64 `flag:"r,reply-rate" help:"回复概率，取值 0~1"`
	Prompt     *string  `flag:"p,prompt" help:"提示词名称"`
	MaxContext *int     `flag:"c,max-context" help:"上下文消息数"`
	Command    *plugin.Command
}

type aiSessionResetCommand struct {
	_       struct{} `cmd:"ai session-reset" help:"重置会话配置为全局默认" usage:"/ai session-reset [-t target]" example:"/ai session-reset\n/ai session-reset -t chatroom:123@chatroom"`
	Target  string   `flag:"t,target" help:"会话 key；不传则重置当前会话"`
	Command *plugin.Command
}

type aiSessionListCommand struct {
	_ struct{} `cmd:"ai session-list" help:"列出所有自定义会话配置" usage:"/ai session-list" example:"/ai session-list"`
}

type aiHelpCommand struct {
	_ struct{} `cmd:"ai help" help:"显示 AI 插件帮助" usage:"/ai help" example:"/ai help"`
}

func registerCommands(p *AiPlugin) error {
	handlers := []func() error{
		func() error { return plugin.RegisterCommand(p.handleGet) },
		func() error { return plugin.RegisterCommand(p.handleSet) },
		func() error { return plugin.RegisterCommand(p.handlePrompt) },
		func() error { return plugin.RegisterCommand(p.handleClearContext) },
		func() error { return plugin.RegisterCommand(p.handleSessionGet) },
		func() error { return plugin.RegisterCommand(p.handleSessionSet) },
		func() error { return plugin.RegisterCommand(p.handleSessionReset) },
		func() error { return plugin.RegisterCommand(p.handleSessionList) },
		//func() error { return plugin.RegisterCommand(p.handleHelp) },
	}
	for _, register := range handlers {
		if err := register(); err != nil {
			return err
		}
	}
	return nil
}

func (p *AiPlugin) handleGet(aiGetCommand) (string, error) {
	config := p.configSnapshot()
	return strings.Join([]string{
		"AI 配置：",
		fmt.Sprintf("url=%s", emptyDash(config.BaseURL)),
		fmt.Sprintf("api_key=%s", maskSecret(config.APIKey)),
		fmt.Sprintf("model=%s", emptyDash(config.Model)),
		fmt.Sprintf("active_prompt=%s", config.ActivePrompt),
		fmt.Sprintf("prompts=%s", strings.Join(promptNames(config.Prompts), ",")),
		fmt.Sprintf("reply_rate=%.4g", config.ReplyRate),
		fmt.Sprintf("max_context=%d", config.MaxContextMessages),
		fmt.Sprintf("timeout=%d", config.HTTPTimeoutSeconds),
		fmt.Sprintf("prompt=%s", emptyDash(activePromptContent(config))),
	}, "\n"), nil
}

func (p *AiPlugin) handleSet(cmd aiSetCommand) (string, error) {
	p.configMu.Lock()
	defer p.configMu.Unlock()

	if err := p.applySetCommandLocked(cmd); err != nil {
		return "", err
	}
	if err := p.SaveConfig(p); err != nil {
		return "", fmt.Errorf("保存 AI 配置失败: %w", err)
	}
	return "AI 配置已更新", nil
}

func (p *AiPlugin) applySetCommandLocked(cmd aiSetCommand) error {
	changed := false
	if cmd.BaseURL != nil {
		p.Config.BaseURL = strings.TrimSpace(*cmd.BaseURL)
		changed = true
	}
	if cmd.APIKey != nil {
		p.Config.APIKey = strings.TrimSpace(*cmd.APIKey)
		changed = true
	}
	if cmd.Model != nil {
		p.Config.Model = strings.TrimSpace(*cmd.Model)
		changed = true
	}
	if cmd.ReplyRate != nil {
		if *cmd.ReplyRate < 0 || *cmd.ReplyRate > 1 {
			return fmt.Errorf("reply_rate 必须在 0~1 之间")
		}
		p.Config.ReplyRate = *cmd.ReplyRate
		changed = true
	}
	if cmd.MaxContext != nil {
		if *cmd.MaxContext <= 0 {
			return fmt.Errorf("max_context 必须大于 0")
		}
		p.Config.MaxContextMessages = *cmd.MaxContext
		changed = true
	}
	if cmd.Timeout != nil {
		if *cmd.Timeout <= 0 {
			return fmt.Errorf("timeout 必须大于 0")
		}
		p.Config.HTTPTimeoutSeconds = *cmd.Timeout
		changed = true
	}
	if !changed {
		return fmt.Errorf("未提供要更新的配置")
	}
	p.Config = normalizeConfigValue(p.Config)
	return nil
}

func (p *AiPlugin) handlePrompt(cmd aiPromptCommand) (string, error) {
	p.configMu.Lock()
	defer p.configMu.Unlock()

	updated, err := p.applyPromptCommandLocked(cmd)
	if err != nil {
		return "", err
	}

	p.clearContext(commandSessionKey(cmd.Command))

	if err := p.SaveConfig(p); err != nil {
		return "", fmt.Errorf("保存 AI 配置失败: %w", err)
	}
	if updated {
		return fmt.Sprintf("提示词已保存并切换：%s", p.Config.ActivePrompt), nil
	}
	return fmt.Sprintf("已切换提示词：%s", p.Config.ActivePrompt), nil
}

func (p *AiPlugin) applyPromptCommandLocked(cmd aiPromptCommand) (bool, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return false, fmt.Errorf("提示词名称不能为空")
	}
	p.Config = normalizeConfigValue(p.Config)
	if cmd.Prompt == nil {
		if _, ok := p.Config.Prompts[name]; !ok {
			return false, fmt.Errorf("提示词不存在：%s", name)
		}
		p.Config.ActivePrompt = name
		return false, nil
	}
	p.Config.Prompts[name] = strings.TrimSpace(*cmd.Prompt)
	p.Config.ActivePrompt = name
	p.Config = normalizeConfigValue(p.Config)
	return true, nil
}

func (p *AiPlugin) handleClearContext(cmd aiClearContextCommand) (string, error) {
	key := strings.TrimSpace(cmd.Target)
	if key == "" {
		key = commandSessionKey(cmd.Command)
	}
	if key == "" {
		return "", fmt.Errorf("无法确定要清理的会话")
	}
	p.clearContext(key)
	return fmt.Sprintf("上下文已清理：%s", key), nil
}

func (p *AiPlugin) handleSessionGet(cmd aiSessionGetCommand) (string, error) {
	key := strings.TrimSpace(cmd.Target)
	if key == "" {
		key = commandSessionKey(cmd.Command)
	}
	if key == "" {
		return "", fmt.Errorf("无法确定要查看的会话")
	}

	config := p.configSnapshot()
	sessionCfg := p.getSessionConfig(key)

	lines := []string{fmt.Sprintf("会话：%s", key)}

	// 回复概率
	if sessionCfg != nil && sessionCfg.ReplyRate != nil {
		lines = append(lines, fmt.Sprintf("reply_rate=%.4g (会话)", *sessionCfg.ReplyRate))
	} else {
		lines = append(lines, fmt.Sprintf("reply_rate=%.4g (全局)", config.ReplyRate))
	}

	// 提示词
	if sessionCfg != nil && sessionCfg.ActivePrompt != nil {
		lines = append(lines, fmt.Sprintf("active_prompt=%s (会话)", *sessionCfg.ActivePrompt))
	} else {
		lines = append(lines, fmt.Sprintf("active_prompt=%s (全局)", config.ActivePrompt))
	}

	// 上下文长度
	if sessionCfg != nil && sessionCfg.MaxContextMessages != nil {
		lines = append(lines, fmt.Sprintf("max_context=%d (会话)", *sessionCfg.MaxContextMessages))
	} else {
		lines = append(lines, fmt.Sprintf("max_context=%d (全局)", config.MaxContextMessages))
	}

	// 显示全局配置
	lines = append(lines, "---")
	lines = append(lines, fmt.Sprintf("全局 reply_rate=%.4g", config.ReplyRate))
	lines = append(lines, fmt.Sprintf("全局 active_prompt=%s", config.ActivePrompt))
	lines = append(lines, fmt.Sprintf("全局 max_context=%d", config.MaxContextMessages))

	return strings.Join(lines, "\n"), nil
}

func (p *AiPlugin) handleSessionSet(cmd aiSessionSetCommand) (string, error) {
	key := strings.TrimSpace(cmd.Target)
	if key == "" {
		key = commandSessionKey(cmd.Command)
	}
	if key == "" {
		return "", fmt.Errorf("无法确定要设置的会话")
	}

	if cmd.ReplyRate == nil && cmd.Prompt == nil && cmd.MaxContext == nil {
		return "", fmt.Errorf("未提供要更新的配置")
	}

	// 验证参数
	if cmd.ReplyRate != nil && (*cmd.ReplyRate < 0 || *cmd.ReplyRate > 1) {
		return "", fmt.Errorf("reply_rate 必须在 0~1 之间")
	}
	if cmd.MaxContext != nil && *cmd.MaxContext <= 0 {
		return "", fmt.Errorf("max_context 必须大于 0")
	}
	if cmd.Prompt != nil {
		config := p.configSnapshot()
		if _, ok := config.Prompts[*cmd.Prompt]; !ok {
			return "", fmt.Errorf("提示词不存在：%s", *cmd.Prompt)
		}
	}

	// 获取或创建会话配置
	p.configMu.Lock()
	defer p.configMu.Unlock()

	sessionCfg := p.getSessionConfig(key)
	if sessionCfg == nil {
		sessionCfg = &SessionConfig{}
	}

	// 更新字段
	if cmd.ReplyRate != nil {
		rate := *cmd.ReplyRate
		sessionCfg.ReplyRate = &rate
	}
	if cmd.Prompt != nil {
		prompt := strings.TrimSpace(*cmd.Prompt)
		sessionCfg.ActivePrompt = &prompt
	}
	if cmd.MaxContext != nil {
		ctx := *cmd.MaxContext
		sessionCfg.MaxContextMessages = &ctx
	}

	p.setSessionConfig(key, sessionCfg)

	if err := p.SaveConfig(p); err != nil {
		return "", fmt.Errorf("保存配置失败: %w", err)
	}

	return fmt.Sprintf("会话配置已更新：%s", key), nil
}

func (p *AiPlugin) handleSessionReset(cmd aiSessionResetCommand) (string, error) {
	key := strings.TrimSpace(cmd.Target)
	if key == "" {
		key = commandSessionKey(cmd.Command)
	}
	if key == "" {
		return "", fmt.Errorf("无法确定要重置的会话")
	}

	p.configMu.Lock()
	defer p.configMu.Unlock()

	p.deleteSessionConfig(key)

	if err := p.SaveConfig(p); err != nil {
		return "", fmt.Errorf("保存配置失败: %w", err)
	}

	return fmt.Sprintf("会话配置已重置为全局默认：%s", key), nil
}

func (p *AiPlugin) handleSessionList(aiSessionListCommand) (string, error) {
	keys := p.listSessionKeys()
	if len(keys) == 0 {
		return "暂无自定义会话配置", nil
	}

	sort.Strings(keys)
	lines := []string{fmt.Sprintf("共 %d 个自定义会话配置：", len(keys))}

	for _, key := range keys {
		cfg := p.getSessionConfig(key)
		if cfg == nil {
			continue
		}

		parts := []string{key}
		if cfg.ReplyRate != nil {
			parts = append(parts, fmt.Sprintf("rate=%.4g", *cfg.ReplyRate))
		}
		if cfg.ActivePrompt != nil {
			parts = append(parts, fmt.Sprintf("prompt=%s", *cfg.ActivePrompt))
		}
		if cfg.MaxContextMessages != nil {
			parts = append(parts, fmt.Sprintf("ctx=%d", *cfg.MaxContextMessages))
		}

		lines = append(lines, strings.Join(parts, ", "))
	}

	return strings.Join(lines, "\n"), nil
}

func commandSessionKey(cmd *plugin.Command) string {
	if cmd == nil || cmd.GetSender() == nil {
		return ""
	}
	sender := cmd.GetSender()
	if sender.GetUsername() == "" {
		return ""
	}
	if sender.GetType() == contactTypeChatroom {
		return "chatroom:" + sender.GetUsername()
	}
	return "private:" + sender.GetUsername()
}

func maskSecret(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "-"
	}
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

func emptyDash(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func promptNames(prompts map[string]string) []string {
	names := make([]string, 0, len(prompts))
	for name := range prompts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
