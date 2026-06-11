package main

import (
	"strings"

	"github.com/sbgayhub/golem/sdk/message"
)

func collectMentionTargets(rule Rule, parsed incomingText, quote quoteInfo, msg *message.Message) []mentionTarget {
	if !rule.At {
		return nil
	}

	targets := make([]mentionTarget, 0, 1)
	targets = appendMentionTarget(targets, mentionTarget{
		Username:    quote.Username,
		DisplayName: quote.Quoter,
	})
	return appendParamMentionTargets(targets, parsed.Param, msg)
}

func appendParamMentionTargets(targets []mentionTarget, param string, msg *message.Message) []mentionTarget {
	names := mentionNamesFromParam(param)
	if len(names) == 0 || msg == nil || msg.GetText() == nil {
		return targets
	}

	reminds := msg.GetText().GetReminds()
	for i := 0; i < len(names) && i < len(reminds); i++ {
		targets = appendMentionTarget(targets, mentionTarget{
			Username:    reminds[i],
			DisplayName: names[i],
		})
	}
	return targets
}

func mentionNamesFromParam(param string) []string {
	tokens := strings.Fields(strings.TrimSpace(param))
	names := make([]string, 0, len(tokens))
	for _, token := range tokens {
		candidate := token
		if _, value, ok := parseNamedArg(token); ok {
			candidate = value
		}
		if !strings.HasPrefix(candidate, "@") {
			continue
		}
		name := strings.TrimSpace(strings.TrimPrefix(candidate, "@"))
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

func appendMentionTarget(targets []mentionTarget, target mentionTarget) []mentionTarget {
	target.Username = strings.TrimSpace(target.Username)
	target.DisplayName = strings.TrimSpace(strings.TrimPrefix(target.DisplayName, "@"))
	if target.Username == "" || target.DisplayName == "" {
		return targets
	}

	for _, existing := range targets {
		if existing.Username == target.Username {
			return targets
		}
	}
	return append(targets, target)
}
