package main

import "strings"

const (
	defaultPromptName         = "default"
	defaultMaxContextMessages = 200
	defaultHTTPTimeoutSeconds = 60
)

// defaultConfig 返回默认配置
func defaultConfig() Config {
	return Config{
		ActivePrompt:       defaultPromptName,
		Prompts:            map[string]string{defaultPromptName: ""},
		ReplyRate:          0.1,
		MaxContextMessages: defaultMaxContextMessages,
		HTTPTimeoutSeconds: defaultHTTPTimeoutSeconds,
	}
}

// normalizeConfig 标准化配置
func (p *AiPlugin) normalizeConfig() {
	p.configMu.Lock()
	defer p.configMu.Unlock()
	p.Config = normalizeConfigValue(p.Config)
}

// configSnapshot 获取配置快照
func (p *AiPlugin) configSnapshot() Config {
	p.configMu.RLock()
	config := p.Config
	p.configMu.RUnlock()
	return normalizeConfigValue(config)
}

// normalizeConfigValue 标准化配置值
func normalizeConfigValue(config Config) Config {
	config.ActivePrompt = strings.TrimSpace(config.ActivePrompt)
	if config.ActivePrompt == "" {
		config.ActivePrompt = defaultPromptName
	}
	config.Prompts = normalizePrompts(config.Prompts, config.LegacyPrompt, config.ActivePrompt)
	config.LegacyPrompt = ""
	if config.MaxContextMessages <= 0 {
		config.MaxContextMessages = defaultMaxContextMessages
	}
	if config.HTTPTimeoutSeconds <= 0 {
		config.HTTPTimeoutSeconds = defaultHTTPTimeoutSeconds
	}
	if config.ReplyRate < 0 {
		config.ReplyRate = 0
	}
	if config.ReplyRate > 1 {
		config.ReplyRate = 1
	}
	config.BaseURL = strings.TrimSpace(config.BaseURL)
	config.APIKey = strings.TrimSpace(config.APIKey)
	config.Model = strings.TrimSpace(config.Model)
	return config
}

// normalizePrompts 标准化提示词映射
func normalizePrompts(prompts map[string]string, legacyPrompt, activePrompt string) map[string]string {
	normalized := make(map[string]string, len(prompts)+1)
	for name, prompt := range prompts {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		normalized[name] = strings.TrimSpace(prompt)
	}
	if legacyPrompt = strings.TrimSpace(legacyPrompt); legacyPrompt != "" {
		if _, ok := normalized[defaultPromptName]; !ok {
			normalized[defaultPromptName] = legacyPrompt
		}
	}
	if _, ok := normalized[activePrompt]; !ok {
		normalized[activePrompt] = ""
	}
	return normalized
}

// activePromptContent 获取活动提示词内容
func activePromptContent(config Config) string {
	return strings.TrimSpace(config.Prompts[config.ActivePrompt])
}

// replyRate 获取全局回复概率（已废弃，保留用于兼容）
func (p *AiPlugin) replyRate() float64 {
	return p.configSnapshot().ReplyRate
}

// maxContextMessages 获取全局上下文消息数（已废弃，保留用于兼容）
func (p *AiPlugin) maxContextMessages() int {
	return p.configSnapshot().MaxContextMessages
}
