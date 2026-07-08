package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sbgayhub/golem/sdk/chatroom"
	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
)

// aiChatPayload 与 ai 插件 ai.chat 能力约定的请求结构
type aiChatPayload struct {
	System   string        `json:"system"`
	Messages []chatMessage `json:"messages"`
}

// chatMessage 发给 ai.chat 能力的消息结构（与 ai 插件的 openAIMessage 对应）
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// historyMsg 从 statistics 表取出的单条历史发言
type historyMsg struct {
	ID        int64
	Content   string
	Timestamp string
}

// profileRecord 画像记录（持久化于 statistics.db 的 profiles 表）
type profileRecord struct {
	Chatroom  string // 群 wxid；全局画像为空字符串
	Member    string // 成员 wxid
	Profile   string // 画像文本
	LastMsgID int64  // 已处理到的 statistics.id 水位线
	UpdatedAt string
}

// named 同时适配 *contact.Contact 与 *chatroom.Member 的展示名读取
type named interface {
	GetRemark() string
	GetNickname() string
	GetAlias() string
	GetUsername() string
}

func displayNameOf(n named) string {
	for _, v := range []string{n.GetRemark(), n.GetNickname(), n.GetAlias(), n.GetUsername()} {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// handleProfile 人物画像触发入口：识别 issuer / chatroomWxid / @ 的 wxid，调用 runProfile，回复结果。
// 返回 (handled=true, nil) 以消费事件，避免 ai 重复回复。
func (p *StatisticsPlugin) handleProfile(msg *message.Message, name string, global, rebuild bool) (bool, error) {
	if msg.GetSender() == nil {
		return true, fmt.Errorf("无法确定消息来源")
	}

	// 异步生成+发送：冷启动可能涉及数十块 × AI 调用，耗时远超事件分发超时（1 分钟）。
	// 若同步处理，分发器会判定 statistics 超时并继续把事件传给 ai，导致「画像 + ai 回复」同时出现。
	// 这里立即返回 handled=true 消费事件（ai 不再回复），后台完成生成与发送。
	go p.processProfile(msg, name, global, rebuild)
	return true, nil
}

// processProfile 实际的画像生成与发送（在后台 goroutine 中运行）。
func (p *StatisticsPlugin) processProfile(msg *message.Message, name string, global, rebuild bool) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("[statistics] 画像生成崩溃", "err", r)
		}
	}()

	// 事件消息自带 Member（真实发消息的人）；群聊里 Sender 是群，Member 是人
	var issuer named
	if m := msg.GetMember(); m != nil {
		issuer = m
	} else {
		issuer = msg.GetSender()
	}
	chatroomWxid := ""
	if msg.Sender.GetType() == contact.ContactType_CONTACT_TYPE_CHATROOM {
		chatroomWxid = msg.GetSender().GetUsername()
	}

	// @ 提人时，微信附带被 @ 人的 wxid（atuserlist/Reminds），优先用于定位
	selfWxid := ""
	if p.contact != nil {
		if self := p.contact.GetSelf(); self != nil {
			selfWxid = self.GetUsername()
		}
	}
	atWxid := extractAtTargetWxid(msg, selfWxid)

	text, err := p.runProfile(issuer, chatroomWxid, name, atWxid, global, rebuild)
	if err != nil {
		slog.Warn("[statistics] 画像生成失败", "err", err)
		errMsg := "画像生成失败：" + err.Error()
		if sendErr := p.sendText(msg.GetSender(), errMsg, nil); sendErr != nil {
			slog.Warn("[statistics] 发送画像错误信息失败", "err", sendErr)
		}
		return
	}
	if text != "" {
		if sendErr := p.sendProfileResult(msg.GetSender(), text); sendErr != nil {
			slog.Warn("[statistics] 发送画像失败", "err", sendErr)
		}
	}
}

// runProfile 先做权限与范围判定，再解析目标成员，最后生成/更新画像。
// issuer 为触发者（群内为真实发言人 Member，私聊为本人）；chatroomWxid 非空表示群聊场景。
// atWxid 为 @ 提人时由微信附带的被 @ 人 wxid（最可靠），非空时优先用于定位目标。
func (p *StatisticsPlugin) runProfile(issuer named, chatroomWxid, targetName, atWxid string, global, rebuild bool) (string, error) {
	if issuer == nil || issuer.GetUsername() == "" {
		return "", fmt.Errorf("无法确定消息来源")
	}
	isChatroom := chatroomWxid != ""
	var owner *contact.Contact
	if p.contact != nil {
		owner = p.contact.GetOwner()
	}

	// 私聊仅主人可查询：非群聊且发信人不是主人则拦截
	if !isChatroom {
		if owner == nil {
			return "私聊查询需先配置主人（Owner）", nil
		}
		if issuer.GetUsername() != owner.GetUsername() {
			return "私聊仅主人可查询人物画像", nil
		}
	}

	var scopeChatroom string // "" => 全局（跨群）
	var memberWxid string
	var displayName string

	switch {
	case global:
		// @ 提人时直接用 wxid；否则按名字在全局联系人缓存中解析
		if atWxid != "" {
			memberWxid = atWxid
			displayName = strings.TrimSpace(targetName)
			if p.contact != nil {
				if c := p.contact.Get(atWxid); c != nil && c.GetUsername() != "" {
					displayName = displayNameOf(c)
				}
			}
		} else {
			if strings.TrimSpace(targetName) == "" {
				return "全局画像需指定成员昵称，例如「人物画像 张三 --global」", nil
			}
			wxid, disp, err := p.resolveGlobal(targetName)
			if err != nil {
				return err.Error(), nil
			}
			memberWxid, displayName = wxid, disp
		}
		scopeChatroom = ""

	case isChatroom:
		// 优先级：@ 提人 wxid > 指定名字 > 未指定则默认查发起人自己
		switch {
		case atWxid != "":
			memberWxid = atWxid
			displayName = strings.TrimSpace(targetName)
			if mem, ok := p.findMemberByWxid(chatroomWxid, atWxid); ok {
				displayName = displayNameOf(mem)
			}
		case strings.TrimSpace(targetName) == "":
			// 群内未指定成员：默认查发起人自己
			memberWxid = issuer.GetUsername()
			displayName = displayNameOf(issuer)
		default:
			mem, ok := p.findMember(chatroomWxid, targetName)
			if !ok {
				return "未在当前群找到成员：" + targetName + "（可尝试 @ 该成员）", nil
			}
			memberWxid = mem.GetUsername()
			displayName = displayNameOf(mem)
		}
		scopeChatroom = chatroomWxid

	default: // 私聊：只能查自己
		if strings.TrimSpace(targetName) != "" && !isSelfName(issuer, targetName) {
			return "私聊中只能查看自己的画像", nil
		}
		memberWxid = issuer.GetUsername()
		displayName = displayNameOf(issuer)
		scopeChatroom = "" // 私聊自己的画像按全局（跨群）处理
	}

	// 主人画像保护：仅主人可查。
	if owner != nil && memberWxid == owner.GetUsername() && issuer.GetUsername() != owner.GetUsername() {
		return "无权查看该成员的画像", nil
	}

	return p.generate(scopeChatroom, memberWxid, displayName, rebuild)
}

// resolveGlobal 按昵称/备注/用户名在全局联系人缓存中解析成员 wxid
func (p *StatisticsPlugin) resolveGlobal(name string) (wxid, display string, err error) {
	if p.contact == nil {
		return "", "", fmt.Errorf("全局未找到成员：%s（contact 能力未注入）", name)
	}
	for _, key := range []string{"nickname::" + name, "remark::" + name, "username::" + name} {
		c := p.contact.Get(key)
		if c != nil && c.GetUsername() != "" {
			return c.GetUsername(), displayNameOf(c), nil
		}
	}
	return "", "", fmt.Errorf("全局未找到成员：%s（非好友无法全局查询，可到其所在群内发「人物画像 %s」）", name, name)
}

// findMember 在当前群成员列表中按群显示名/昵称/备注/用户名匹配。
// 群聊 @ 提人插入的是成员的「群显示名(DisplayName)」，可能与真实昵称不同，
// 因此 DisplayName 必须参与匹配，且优先。
func (p *StatisticsPlugin) findMember(chatroomWxid, name string) (*chatroom.Member, bool) {
	if p.chatroom == nil {
		return nil, false
	}
	name = strings.TrimSpace(name)
	for _, m := range p.chatroom.ListMembers(chatroomWxid) {
		if m == nil {
			continue
		}
		for _, v := range []string{m.GetDisplayName(), m.GetRemark(), m.GetNickname(), m.GetAlias(), m.GetUsername()} {
			if strings.EqualFold(strings.TrimSpace(v), name) {
				return m, true
			}
		}
	}
	return nil, false
}

// findMemberByWxid 在当前群成员列表中按 wxid 精确匹配（用于 @ 提人时已拿到 wxid 的场景）
func (p *StatisticsPlugin) findMemberByWxid(chatroomWxid, wxid string) (*chatroom.Member, bool) {
	if p.chatroom == nil {
		return nil, false
	}
	for _, m := range p.chatroom.ListMembers(chatroomWxid) {
		if m != nil && m.GetUsername() == wxid {
			return m, true
		}
	}
	return nil, false
}

func isSelfName(sender named, name string) bool {
	name = strings.TrimSpace(name)
	for _, v := range []string{sender.GetRemark(), sender.GetNickname(), sender.GetAlias(), sender.GetUsername()} {
		if strings.EqualFold(strings.TrimSpace(v), name) {
			return true
		}
	}
	return false
}

// generate 读取历史发言 → 切块 → 调用 ai.chat → 合并/冷启动 → 持久化
func (p *StatisticsPlugin) generate(scopeChatroom, memberWxid, displayName string, rebuild bool) (string, error) {
	if p.store == nil {
		return "", fmt.Errorf("存储未初始化")
	}
	cfg := normalizeConfig(p.Config)
	rec, exists := p.store.loadProfile(scopeChatroom, memberWxid)

	sinceID := int64(0)
	limit := 0
	if exists && !rebuild {
		sinceID = rec.LastMsgID
	} else {
		limit = cfg.ColdStartMaxMessages // 冷启动安全天花板
	}

	msgs, err := p.store.queryHistory(scopeChatroom, memberWxid, sinceID, limit)
	if err != nil {
		return "", err
	}
	if len(msgs) == 0 {
		if exists && strings.TrimSpace(rec.Profile) != "" {
			// 已有画像且无新发言：直接把已有画像发出来，不空跑一次 AI
			return rec.Profile, nil
		}
		if exists {
			return "该成员暂无新发言，且尚无已生成的画像", nil
		}
		return "该成员暂无可分析的发言记录", nil
	}

	chunks := splitIntoChunks(msgs, cfg.ChunkTokenBudget, cfg.MaxSingleMsgChars)
	if len(chunks) > cfg.ColdStartMaxChunks {
		chunks = keepRecentChunks(chunks, cfg.ColdStartMaxChunks)
	}
	// 水位线必须基于「实际送入模型的块」计算，否则采样丢弃的块会被永久跳过
	coveredID := maxChunkID(chunks)

	observations := make([]string, 0, len(chunks))
	for _, ch := range chunks {
		obs, err := p.callAIChunk(ch)
		if err != nil {
			return "", fmt.Errorf("分析发言片段失败: %w", err)
		}
		observations = append(observations, obs)
	}

	quant := summarizeQuant(msgs)
	var final string
	if exists && !rebuild {
		final, err = p.callAIMerge(displayName, rec.Profile, observations, quant)
	} else {
		final, err = p.callAIMerge(displayName, "", observations, quant)
	}
	if err != nil {
		return "", fmt.Errorf("生成画像失败: %w", err)
	}
	if strings.TrimSpace(final) == "" {
		return "", fmt.Errorf("模型返回空画像（可能上下文不足或调用异常）")
	}

	if err := p.store.saveProfile(profileRecord{
		Chatroom:  scopeChatroom,
		Member:    memberWxid,
		Profile:   final,
		LastMsgID: coveredID,
	}); err != nil {
		return "", fmt.Errorf("保存画像失败: %w", err)
	}
	slog.Info("[statistics] 画像已生成", "scope", scopeChatroom, "member", memberWxid, "chunks", len(chunks))
	return final, nil
}

// callAIChunk 对单块历史发言产出局部观察
func (p *StatisticsPlugin) callAIChunk(ch []historyMsg) (string, error) {
	payload, err := json.Marshal(aiChatPayload{
		System:   systemColdChunk,
		Messages: []chatMessage{{Role: "user", Content: formatChunk(ch)}},
	})
	if err != nil {
		return "", err
	}
	return p.callAI(string(payload))
}

// callAIMerge 合并已有画像与新增观察，产出完整画像
func (p *StatisticsPlugin) callAIMerge(displayName, existing string, observations []string, quant string) (string, error) {
	user := buildMergeUserContent(displayName, existing, observations, quant)
	payload, err := json.Marshal(aiChatPayload{
		System:   systemMerge,
		Messages: []chatMessage{{Role: "user", Content: user}},
	})
	if err != nil {
		return "", err
	}
	return p.callAI(string(payload))
}

// callAI 经跨插件调用使用 ai 插件的 ai.chat 能力
func (p *StatisticsPlugin) callAI(payload string) (string, error) {
	if p.caller == nil {
		return "", fmt.Errorf("调用能力未注入（需要 ai 插件提供 ai.chat）")
	}
	mime, data, err := p.caller.CallPlugin("ai.chat", map[string]string{"payload": payload})
	if err != nil {
		return "", err
	}
	_ = mime
	return strings.TrimSpace(string(data)), nil
}

// summarizeQuant 基于本次拉取的历史发言生成量化指标
func summarizeQuant(msgs []historyMsg) string {
	if len(msgs) == 0 {
		return "无"
	}
	total := 0
	for _, m := range msgs {
		total += len([]rune(m.Content))
	}
	first, last := msgs[0].Timestamp, msgs[len(msgs)-1].Timestamp
	if first == "" {
		first = "未知"
	}
	if last == "" {
		last = "未知"
	}
	return fmt.Sprintf("发言条数: %d\n时间跨度: %s ~ %s\n平均每条字数: %d",
		len(msgs), first, last, total/len(msgs))
}
