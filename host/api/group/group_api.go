// Package groupapi 提供群聊服务的 API 接口定义。
package groupapi

// GroupService 群聊服务 API 接口（返回 API proto 类型）
type GroupService interface {
	// Create 创建群聊
	Create(members []string) (*CreateGroupResponse, error)
	// FacingCreate 面对面建群
	FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateGroupResponse, error)
	// GetInfo 获取群详细信息
	GetInfo(groupID string) (*GetGroupInfoDetailResponse, error)
	// GetMemberDetail 获取群成员详情
	GetMemberDetail(groupID string) (*GetGroupMemberDetailResponse, error)
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
}
