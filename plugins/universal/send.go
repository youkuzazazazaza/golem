package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
)

func (p *UniversalPlugin) sendResult(receiver *contact.Contact, sendType, result string, mentions []mentionTarget) error {
	if p.message == nil {
		return errors.New("message ability is not injected")
	}

	var mediaData []byte
	if needsMediaData(sendType) {
		var err error
		mediaData, err = downloadMedia(result, p.requestTimeout())
		if err != nil {
			return err
		}
	}

	msg, err := buildMessage(receiver, sendType, result, mediaData, mentions)
	if err != nil {
		return err
	}
	_, err = p.message.Send(msg)
	return err
}

func buildMessage(receiver *contact.Contact, sendType, result string, mediaData []byte, mentions []mentionTarget) (*message.Message, error) {
	if receiver == nil {
		return nil, errors.New("receiver is empty")
	}
	if strings.TrimSpace(result) == "" {
		return nil, errors.New("result is empty")
	}

	content := strings.TrimRight(strings.TrimSpace(result), "\n")
	var reminds []string
	if normalizeSendType(sendType) == "text" {
		content, reminds = applyMentionPrefix(result, mentions)
	}
	msg := &message.Message{
		Receiver: receiver,
		Content:  content,
	}
	switch normalizeSendType(sendType) {
	case "text":
		msg.Type = message.TypeText
		msg.Data = &message.Message_Text{Text: &message.TextData{Content: content, Reminds: reminds}}
	case "image":
		if len(mediaData) == 0 {
			return nil, errors.New("image data is empty")
		}
		msg.Type = message.TypeImage
		msg.Data = &message.Message_Image{Image: &message.ImageData{Media: &message.Media{Data: mediaData}}}
	case "video":
		if len(mediaData) == 0 {
			return nil, errors.New("video data is empty")
		}
		msg.Type = message.TypeVideo
		msg.Data = &message.Message_Video{Video: &message.VideoData{Media: &message.Media{Data: mediaData}}}
	case "emoji":
		if len(mediaData) == 0 {
			return nil, errors.New("emoji data is empty")
		}
		msg.Type = message.TypeEmoji
		msg.Data = &message.Message_Emoji{Emoji: &message.EmojiData{Media: &message.Media{Data: mediaData}}}
	default:
		return nil, fmt.Errorf("unsupported send_type: %s", sendType)
	}
	return msg, nil
}

func applyMentionPrefix(result string, mentions []mentionTarget) (string, []string) {
	prefixes := make([]string, 0, len(mentions))
	reminds := make([]string, 0, len(mentions))
	for _, mention := range mentions {
		displayName := strings.TrimSpace(strings.TrimPrefix(mention.DisplayName, "@"))
		username := strings.TrimSpace(mention.Username)
		if displayName == "" || username == "" {
			continue
		}
		prefixes = append(prefixes, "@"+displayName)
		reminds = append(reminds, username)
	}
	if len(prefixes) == 0 {
		return result, nil
	}
	return strings.Join(prefixes, " ") + " " + result, reminds
}

func needsMediaData(sendType string) bool {
	switch normalizeSendType(sendType) {
	case "image", "video", "emoji":
		return true
	default:
		return false
	}
}

func downloadMedia(rawURL string, timeout time.Duration) ([]byte, error) {
	data, err := doRequest(http.MethodGet, strings.TrimSpace(rawURL), nil, "", timeout)
	if err != nil {
		return nil, fmt.Errorf("download media: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("download media: empty body")
	}
	return data, nil
}
