package main

import (
	"strings"
)

func emptyPlaceholder(value string) string {
	if value == "" {
		return "(空)"
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func truncate(data []byte, maxLen int) string {
	text := string(data)
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
