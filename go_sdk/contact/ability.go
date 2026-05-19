package contact

import (
	"github.com/sbgayhub/golem/sdk/group"
)

// Ability 联系人能力接口（供插件嵌入使用）
type Ability interface {
	// GetContactByKey 按键查询缓存联系人
	// 支持前缀：username::（默认）、nickname::、remark::
	GetContactByKey(key string) (*Contact, bool)
	// GetContactByStrategy 按策略查询缓存联系人
	GetContactByStrategy(key string, strategy RetrievalType) (*Contact, bool)
	// GetContactList 获取联系人列表
	GetContactList() ([]*Contact, error)
	// GetGroupMembers 获取群成员列表
	GetGroupMembers(groupId string) ([]*group.GroupMember, bool)
	// SetRemark 设置联系人备注
	SetRemark(username, remark string) (*OperateResponse, error)
	// AddFriend 发送好友申请
	AddFriend(v1, v2, content string, operate, scene int) (*OperateResponse, error)
	// VerifyFriend 通过好友验证
	VerifyFriend(v1, v2 string, scene int) (*OperateResponse, error)
	// Delete 删除联系人
	Delete(username string) (*OperateResponse, error)
	// BlacklistAdd 添加到黑名单
	BlacklistAdd(username string) (*OperateResponse, error)
	// BlacklistRemove 从黑名单移除
	BlacklistRemove(username string) (*OperateResponse, error)
	// Search 搜索联系人
	Search(keyword string, fromScene, searchScene uint32) ([]*Contact, error)
}

// Instance 联系人能力实例（由 host/ability 层注入）
var Instance Ability
