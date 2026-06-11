package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

func renderTemplate(raw string, vars map[string]string) (string, error) {
	if raw == "" {
		return "", nil
	}
	funcs := template.FuncMap{}
	for key, value := range vars {
		if !templateNamePattern.MatchString(key) {
			continue
		}
		funcs[key] = func(v string) func() string {
			return func() string { return v }
		}(value)
	}

	tpl, err := template.New("universal").Funcs(funcs).Parse(raw)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	if err := tpl.Execute(&builder, nil); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func buildTemplateVars(parsed incomingText, quote quoteInfo) map[string]string {
	vars := map[string]string{
		"text":    parsed.Text,
		"keyword": parsed.Keyword,
		"param":   parsed.Param,
		"quote":   quote.Quote,
		"quoter":  quote.Quoter,
	}

	argIndex := 1
	if quote.Quoter != "" {
		vars[fmt.Sprintf("arg_%d", argIndex)] = quote.Quoter
		argIndex++
	}

	for _, token := range strings.Fields(parsed.Param) {
		if name, value, ok := parseNamedArg(token); ok {
			if !isReservedTemplateName(name) {
				vars[name] = strings.TrimPrefix(value, "@")
			}
			continue
		}
		vars[fmt.Sprintf("arg_%d", argIndex)] = strings.TrimPrefix(token, "@")
		argIndex++
	}

	return vars
}

func parseNamedArg(token string) (string, string, bool) {
	name, value, ok := strings.Cut(token, "=")
	if !ok || !templateNamePattern.MatchString(name) {
		return "", "", false
	}
	return name, value, true
}

func isReservedTemplateName(name string) bool {
	switch name {
	case "text", "keyword", "param", "quote", "quoter":
		return true
	}
	if strings.HasPrefix(name, "arg_") {
		_, err := strconv.Atoi(strings.TrimPrefix(name, "arg_"))
		return err == nil
	}
	return false
}
