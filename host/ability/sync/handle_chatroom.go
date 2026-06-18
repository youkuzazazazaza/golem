package sync

import (
	chatroomability "github.com/sbgayhub/golem/host/ability/chatroom"
	messageapi "github.com/sbgayhub/golem/host/api/message"
)

func handleChatroomMember(members []*messageapi.ModifyChatroomMember) {
	publishChangeEvents(chatroomability.HandleModifyMembers(members))
}

func handleChatroomMemberDisplayName(displayNames []*messageapi.ModifyChatroomMemberDisplayName) {
	publishChangeEvents(chatroomability.HandleModifyMemberDisplayNames(displayNames))
}
