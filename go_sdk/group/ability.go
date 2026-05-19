package group

import "time"

// GroupCache 群缓存结构
type GroupCache struct {
	data      map[string]*Group
	updatedAt map[string]time.Time
	ttl       time.Duration // 预留 TTL，暂不使用
}

// NewGroupCache 创建群缓存
func NewGroupCache() *GroupCache {
	return &GroupCache{
		data:      make(map[string]*Group),
		updatedAt: make(map[string]time.Time),
	}
}

// GroupMemberCache 群成员缓存结构
type GroupMemberCache struct {
	data      map[string][]*GroupMember
	updatedAt map[string]time.Time
	ttl       time.Duration // 预留 TTL，暂不使用
}

// NewGroupMemberCache 创建群成员缓存
func NewGroupMemberCache() *GroupMemberCache {
	return &GroupMemberCache{
		data:      make(map[string][]*GroupMember),
		updatedAt: make(map[string]time.Time),
	}
}

// Ability 群聊能力接口（供插件嵌入使用）
type Ability interface {
	// Create 创建群聊
	Create(members []string) (*CreateGroupResponse, error)
	// FacingCreate 面对面建群
	FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateGroupResponse, error)
	// GetInfo 获取群信息
	GetInfo(groupID string) (*Group, error)
	// GetMemberDetail 获取群成员详情
	GetMemberDetail(groupID string) ([]*GroupMember, error)
	// GetQRCode 获取群二维码
	GetQRCode(groupID string) (*GetGroupQRCodeResponse, error)
	// AddMember 添加群成员
	AddMember(groupID string, members []string) (*AddGroupMemberResponse, error)
	// InviteMember 邀请群成员
	InviteMember(groupID string, members []string) (*InviteGroupMemberResponse, error)
	// RemoveMember 移除群成员
	RemoveMember(groupID string, members []string) (*RemoveGroupMemberResponse, error)
	// SetName 设置群名称
	SetName(groupID, name string) (*OperateResponse, error)
	// SetAnnouncement 设置群公告
	SetAnnouncement(groupID, content string) (*SetAnnouncementResponse, error)
	// SetRemark 设置群备注
	SetRemark(groupID, remark string) (*OperateResponse, error)
	// SetContactList 保存到通讯录
	SetContactList(groupID string, save bool) (*OperateResponse, error)
	// SetAdmin 设置群管理员
	SetAdmin(groupID string, members []string) (*OperateResponse, error)
	// RemoveAdmin 移除群管理员
	RemoveAdmin(groupID string, members []string) (*OperateResponse, error)
	// TransferOwner 转让群主
	TransferOwner(groupID, newOwner string) (*OperateResponse, error)
	// Quit 退出群聊
	Quit(groupID string) (*OperateResponse, error)
	// ScanJoin 扫码进群
	ScanJoin(qrcodeURL string) (*JoinResult, error)
	// ScanJoinEnterprise 企业微信扫码进群
	ScanJoinEnterprise(qrcodeURL string) (*JoinResult, error)
	// ConsentJoin 同意入群邀请
	ConsentJoin(inviteURL string) (*JoinResult, error)
	// GetGroupByKey 按键查询缓存群信息
	// 支持前缀：group_id::（默认）、name::
	GetGroupByKey(key string) (*Group, bool)
	// GetGroupByStrategy 按策略查询缓存群信息
	GetGroupByStrategy(key string, strategy RetrievalType) (*Group, bool)
	// GetGroupMembers 获取缓存群成员列表
	GetGroupMembers(groupID string) ([]*GroupMember, bool)
}

// Instance 群聊能力实例（由 host/ability 层注入）
var Instance Ability
