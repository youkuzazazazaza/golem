package sync

import (
	"fmt"
	"log/slog"

	messageability "github.com/sbgayhub/golem/host/ability/message"
	messageapi "github.com/sbgayhub/golem/host/api/message"
	contactsdk "github.com/sbgayhub/golem/sdk/contact"
	messagesdk "github.com/sbgayhub/golem/sdk/message"
)

func handleMessage(messages []*messageapi.NewMessage) {
	for _, msg := range messages {
		if data, err := messageability.Build(msg); err != nil {
			slog.Error("构建消息失败", "err", err)
		} else {
			log(data)

			//// 检查是否为命令
			//if cmd, ok := plugin.ParseCommand(data.GetBase().Content, data.GetBase().Sender); ok {
			//	slog.Info("收到命令", "main", cmd.Main, "sub", cmd.Sub, "args", cmd.Args)
			//	plugin.DispatchCommand(cmd)
			//} else {
			//	// 普通消息，发布事件
			//	plugin.Publish(&pluginsdk.Event{
			//		Topic:   data.GetType().Topic,
			//		Sender:  data.GetBase().Sender,
			//		Payload: data,
			//	})
			//}
		}
	}
}

func log(message *messagesdk.Message) {
	if message.Sender.GetType() == contactsdk.ContactType_CONTACT_TYPE_GROUP {
		sender := "system"
		if message.Member != nil {
			sender = message.Member.Nickname
		}
		slog.Info(fmt.Sprintf("%s -> %s: [%s] %s", sender, message.Sender.GetNickname(), message.Type.Desc, message.Content))
	} else {
		slog.Info(fmt.Sprintf("%s -> %s: [%s] %s", message.Sender.GetNickname(), message.Receiver.GetNickname(), message.Type.Desc, message.Content))
	}
}
