package sync

import (
	"log/slog"

	contactability "github.com/sbgayhub/golem/host/ability/contact"
	messageapi "github.com/sbgayhub/golem/host/api/message"
	pluginsdk "github.com/sbgayhub/golem/sdk/plugin"
)

func handleUserInfo(infos []*messageapi.ModifyUserInfo, extends []*messageapi.UserInfoExtend) {
	maxLen := len(infos)
	if len(extends) > maxLen {
		maxLen = len(extends)
	}
	for index := 0; index < maxLen; index++ {
		var info *messageapi.ModifyUserInfo
		var ext *messageapi.UserInfoExtend
		if index < len(infos) {
			info = infos[index]
		}
		if index < len(extends) {
			ext = extends[index]
		}
		event, err := contactability.ApplySelfInfo(info, ext)
		if err != nil {
			slog.Warn("处理用户信息变更失败", "err", err)
			continue
		}
		if event != nil {
			publishChangeEvents([]*pluginsdk.ChangeEvent{event})
		}
	}
}
