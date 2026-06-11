package main

import (
	"encoding/xml"
	"strings"

	"github.com/sbgayhub/golem/sdk/message"
	"github.com/tidwall/gjson"
)

func parseIncomingText(text string) (incomingText, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return incomingText{}, false
	}

	keyword, param, hasParam := strings.Cut(text, " ")
	if !hasParam {
		return incomingText{Text: text, Keyword: keyword}, true
	}
	return incomingText{
		Text:    text,
		Keyword: strings.TrimSpace(keyword),
		Param:   strings.TrimSpace(param),
	}, keyword != ""
}

func extractQuote(msg *message.Message) quoteInfo {
	if msg == nil {
		return quoteInfo{}
	}
	if app := msg.GetApp(); app != nil {
		if quote := parseQuoteXML(app.GetXml()); quote.Quoter != "" || quote.Quote != "" {
			return quote
		}
	}
	if raw := msg.GetRaw(); raw != "" {
		if content := gjson.Get(raw, "content.value"); content.Exists() {
			if quote := parseQuoteXML(content.String()); quote.Quoter != "" || quote.Quote != "" {
				return quote
			}
		}
	}
	return quoteInfo{}
}

func parseQuoteXML(raw string) quoteInfo {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return quoteInfo{}
	}

	var data struct {
		AppMsg struct {
			Refer quoteRefer `xml:"refermsg"`
		} `xml:"appmsg"`
		Refer quoteRefer `xml:"refermsg"`
	}
	if err := xml.Unmarshal([]byte(raw), &data); err != nil {
		return quoteInfo{}
	}
	refer := data.AppMsg.Refer
	if refer.Content == "" && refer.DisplayName == "" && refer.FromUser == "" && refer.ChatUser == "" {
		refer = data.Refer
	}
	username := strings.TrimSpace(refer.ChatUser)
	return quoteInfo{
		Username: username,
		Quoter:   firstNonEmpty(refer.DisplayName, username, refer.FromUser),
		Quote:    strings.TrimSpace(refer.Content),
	}
}

func messageContent(msg *message.Message) string {
	if msg == nil {
		return ""
	}
	if text := msg.GetText(); text != nil && text.GetContent() != "" {
		return text.GetContent()
	}
	return msg.GetContent()
}
