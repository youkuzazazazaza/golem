package main

import (
	"strconv"
	"strings"
)

// systemColdChunk 冷启动：对单块历史发言产出“局部观察”，不下总结性结论。
const systemColdChunk = `你正在分析一位群成员的历史发言片段。请只做客观观察、不要下总结性结论，用中文要点列表输出，每条简短。重点提取：
- 语言风格（口语化/书面/表情包/标点习惯）
- 常聊话题（游戏、美食、工作、八卦等）
- 活跃时间特征（从发言内容推断，如夜猫子）
- 口头禅、高频词、常用表情
- 与他人互动的方式（怼人、捧场、潜水等）
若片段信息不足，如实说明“信息不足”。`

// systemMergeForRender 合并阶段系统提示：结合已有画像与新增观察，产出完整人物画像。
// renderImage 在源头控制 LLM 输出格式，免去渲染前剥离 markdown 的工作：
//   - true：要求输出 markdown（## 标题、**加粗**、- 列表），供 markdown.to.image 渲染
//   - false：要求输出纯文本（无任何 markdown 符号），供直接发送
func systemMergeForRender(renderImage bool) string {
	base := `你正在维护一位群成员的人物画像。下面会给出【已有画像】和【新增发言观察】以及【量化指标】。
请合并更新，输出一份完整、连贯的人物画像，用中文、条理清晰。维度包括：
1. 性格与语言风格
2. 兴趣话题
3. 活跃时段
4. 口头禅 / 表情习惯
5. 与其他成员的关系（互动对象、关系亲疏）
可参考量化指标（发言量、时间跨度等），但画像要基于实际发言内容，不要编造。
若【已有画像】为空，则直接基于【新增发言观察】生成画像。`
	if renderImage {
		return base + "\n\n输出格式要求：请用 markdown 格式输出，使用 ## 标题、**加粗**、- 列表等结构增强可读性，以便后续渲染为图片。"
	}
	return base + "\n\n输出格式要求：请直接输出纯文本，不要使用任何 markdown 符号（如 #、**、*、-、>、` 等），用换行和空行组织结构，以便直接发送文本。"
}

// buildMergeUserContent 组装合并阶段的 user 消息
func buildMergeUserContent(displayName, existing string, observations []string, quant string) string {
	out := "成员昵称：" + displayName + "\n\n"
	out += "【已有画像】\n"
	if strings.TrimSpace(existing) == "" {
		out += "(无)\n"
	} else {
		out += existing + "\n"
	}
	out += "\n【新增发言观察】\n"
	for i, o := range observations {
		out += "· 片段" + strconv.Itoa(i+1) + "：\n" + o + "\n"
	}
	out += "\n【量化指标】\n" + quant + "\n"
	return out
}
