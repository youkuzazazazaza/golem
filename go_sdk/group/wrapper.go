package group

import (
	"context"
)

// Client 实现 Ability 接口，通过 gRPC 调用远程群聊服务
type Client struct {
	Client GroupServiceClient
}

var _ Ability = (*Client)(nil)

// Create 创建群聊
func (c Client) Create(members []string) (*CreateGroupResponse, error) {
	return c.Client.Create(context.Background(), &CreateGroupRequest{Members: members})
}

// FacingCreate 面对面建群
func (c Client) FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateGroupResponse, error) {
	return c.Client.FacingCreate(context.Background(), &FacingCreateGroupRequest{
		Password:  password,
		Latitude:  latitude,
		Longitude: longitude,
		Operate:   operate,
	})
}

// GetInfo 获取群信息
func (c Client) GetInfo(groupID string) (*Group, error) {
	resp, err := c.Client.GetInfo(context.Background(), &GetGroupInfoRequest{GroupId: groupID})
	if err != nil {
		return nil, err
	}
	return resp.Info, nil
}

// GetMemberDetail 获取群成员详情
func (c Client) GetMemberDetail(groupID string) ([]*GroupMember, error) {
	resp, err := c.Client.GetMemberDetail(context.Background(), &GetGroupMemberDetailRequest{GroupId: groupID})
	if err != nil {
		return nil, err
	}
	return resp.Members, nil
}

// GetQRCode 获取群二维码
func (c Client) GetQRCode(groupID string) (*GetGroupQRCodeResponse, error) {
	return c.Client.GetQRCode(context.Background(), &GetGroupQRCodeRequest{GroupId: groupID})
}

// AddMember 添加群成员
func (c Client) AddMember(groupID string, members []string) (*AddGroupMemberResponse, error) {
	return c.Client.AddMember(context.Background(), &AddGroupMemberRequest{
		GroupId: groupID,
		Members: members,
	})
}

// InviteMember 邀请群成员
func (c Client) InviteMember(groupID string, members []string) (*InviteGroupMemberResponse, error) {
	return c.Client.InviteMember(context.Background(), &InviteGroupMemberRequest{
		GroupId: groupID,
		Members: members,
	})
}

// RemoveMember 移除群成员
func (c Client) RemoveMember(groupID string, members []string) (*RemoveGroupMemberResponse, error) {
	return c.Client.RemoveMember(context.Background(), &RemoveGroupMemberRequest{
		GroupId: groupID,
		Members: members,
	})
}

// SetName 设置群名称
func (c Client) SetName(groupID, name string) (*OperateResponse, error) {
	return c.Client.SetName(context.Background(), &SetGroupNameRequest{
		GroupId: groupID,
		Name:    name,
	})
}

// SetAnnouncement 设置群公告
func (c Client) SetAnnouncement(groupID, content string) (*SetAnnouncementResponse, error) {
	return c.Client.SetAnnouncement(context.Background(), &SetAnnouncementRequest{
		GroupId: groupID,
		Content: content,
	})
}

// SetRemark 设置群备注
func (c Client) SetRemark(groupID, remark string) (*OperateResponse, error) {
	return c.Client.SetRemark(context.Background(), &SetGroupRemarkRequest{
		GroupId: groupID,
		Remark:  remark,
	})
}

// SetContactList 保存到通讯录
func (c Client) SetContactList(groupID string, save bool) (*OperateResponse, error) {
	return c.Client.SetContactList(context.Background(), &SetContactListRequest{
		GroupId: groupID,
		Save:    save,
	})
}

// SetAdmin 设置群管理员
func (c Client) SetAdmin(groupID string, members []string) (*OperateResponse, error) {
	return c.Client.SetAdmin(context.Background(), &SetGroupAdminRequest{
		GroupId: groupID,
		Members: members,
	})
}

// RemoveAdmin 移除群管理员
func (c Client) RemoveAdmin(groupID string, members []string) (*OperateResponse, error) {
	return c.Client.RemoveAdmin(context.Background(), &RemoveGroupAdminRequest{
		GroupId: groupID,
		Members: members,
	})
}

// TransferOwner 转让群主
func (c Client) TransferOwner(groupID, newOwner string) (*OperateResponse, error) {
	return c.Client.TransferOwner(context.Background(), &TransferGroupOwnerRequest{
		GroupId:  groupID,
		NewOwner: newOwner,
	})
}

// Quit 退出群聊
func (c Client) Quit(groupID string) (*OperateResponse, error) {
	return c.Client.Quit(context.Background(), &QuitGroupRequest{GroupId: groupID})
}

// ScanJoin 扫码进群
func (c Client) ScanJoin(qrcodeURL string) (*JoinResult, error) {
	return c.Client.ScanJoin(context.Background(), &ScanJoinRequest{QrcodeUrl: qrcodeURL})
}

// ScanJoinEnterprise 企业微信扫码进群
func (c Client) ScanJoinEnterprise(qrcodeURL string) (*JoinResult, error) {
	return c.Client.ScanJoinEnterprise(context.Background(), &ScanJoinRequest{QrcodeUrl: qrcodeURL})
}

// ConsentJoin 同意入群邀请
func (c Client) ConsentJoin(inviteURL string) (*JoinResult, error) {
	return c.Client.ConsentJoin(context.Background(), &ConsentJoinRequest{InviteUrl: inviteURL})
}

// GetGroupByKey 按键查询缓存群信息
func (c Client) GetGroupByKey(key string) (*Group, bool) {
	resp, err := c.Client.GetGroupByKey(context.Background(), &GetGroupByKeyRequest{Key: key})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Group, true
}

// GetGroupByStrategy 按策略查询缓存群信息
func (c Client) GetGroupByStrategy(key string, strategy RetrievalType) (*Group, bool) {
	resp, err := c.Client.GetGroupByStrategy(context.Background(), &GetGroupByStrategyRequest{
		Key:      key,
		Strategy: strategy,
	})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Group, true
}

// GetGroupMembers 获取缓存群成员列表
func (c Client) GetGroupMembers(groupID string) ([]*GroupMember, bool) {
	resp, err := c.Client.GetGroupMembers(context.Background(), &GetGroupMembersByKeyRequest{GroupId: groupID})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Members, true
}

// Server 实现 GroupServiceServer 接口，将 gRPC 请求委托给 Ability 实现
type Server struct {
	UnimplementedGroupServiceServer
	Impl Ability
}

// Create 创建群聊
func (s Server) Create(ctx context.Context, request *CreateGroupRequest) (*CreateGroupResponse, error) {
	return s.Impl.Create(request.Members)
}

// FacingCreate 面对面建群
func (s Server) FacingCreate(ctx context.Context, request *FacingCreateGroupRequest) (*FacingCreateGroupResponse, error) {
	return s.Impl.FacingCreate(request.Password, request.Latitude, request.Longitude, request.Operate)
}

// GetInfo 获取群信息
func (s Server) GetInfo(ctx context.Context, request *GetGroupInfoRequest) (*GetGroupResponse, error) {
	info, err := s.Impl.GetInfo(request.GroupId)
	if err != nil {
		return nil, err
	}
	return &GetGroupResponse{Group: info}, nil
}

// GetMemberDetail 获取群成员详情
func (s Server) GetMemberDetail(ctx context.Context, request *GetGroupMemberDetailRequest) (*GetGroupMemberDetailResponse, error) {
	members, err := s.Impl.GetMemberDetail(request.GroupId)
	if err != nil {
		return nil, err
	}
	return &GetGroupMemberDetailResponse{GroupId: request.GroupId, Members: members}, nil
}

// GetQRCode 获取群二维码
func (s Server) GetQRCode(ctx context.Context, request *GetGroupQRCodeRequest) (*GetGroupQRCodeResponse, error) {
	return s.Impl.GetQRCode(request.GroupId)
}

// AddMember 添加群成员
func (s Server) AddMember(ctx context.Context, request *AddGroupMemberRequest) (*AddGroupMemberResponse, error) {
	return s.Impl.AddMember(request.GroupId, request.Members)
}

// InviteMember 邀请群成员
func (s Server) InviteMember(ctx context.Context, request *InviteGroupMemberRequest) (*InviteGroupMemberResponse, error) {
	return s.Impl.InviteMember(request.GroupId, request.Members)
}

// RemoveMember 移除群成员
func (s Server) RemoveMember(ctx context.Context, request *RemoveGroupMemberRequest) (*RemoveGroupMemberResponse, error) {
	return s.Impl.RemoveMember(request.GroupId, request.Members)
}

// SetName 设置群名称
func (s Server) SetName(ctx context.Context, request *SetGroupNameRequest) (*OperateResponse, error) {
	return s.Impl.SetName(request.GroupId, request.Name)
}

// SetAnnouncement 设置群公告
func (s Server) SetAnnouncement(ctx context.Context, request *SetAnnouncementRequest) (*SetAnnouncementResponse, error) {
	return s.Impl.SetAnnouncement(request.GroupId, request.Content)
}

// SetRemark 设置群备注
func (s Server) SetRemark(ctx context.Context, request *SetGroupRemarkRequest) (*OperateResponse, error) {
	return s.Impl.SetRemark(request.GroupId, request.Remark)
}

// SetContactList 保存到通讯录
func (s Server) SetContactList(ctx context.Context, request *SetContactListRequest) (*OperateResponse, error) {
	return s.Impl.SetContactList(request.GroupId, request.Save)
}

// SetAdmin 设置群管理员
func (s Server) SetAdmin(ctx context.Context, request *SetGroupAdminRequest) (*OperateResponse, error) {
	return s.Impl.SetAdmin(request.GroupId, request.Members)
}

// RemoveAdmin 移除群管理员
func (s Server) RemoveAdmin(ctx context.Context, request *RemoveGroupAdminRequest) (*OperateResponse, error) {
	return s.Impl.RemoveAdmin(request.GroupId, request.Members)
}

// TransferOwner 转让群主
func (s Server) TransferOwner(ctx context.Context, request *TransferGroupOwnerRequest) (*OperateResponse, error) {
	return s.Impl.TransferOwner(request.GroupId, request.NewOwner)
}

// Quit 退出群聊
func (s Server) Quit(ctx context.Context, request *QuitGroupRequest) (*OperateResponse, error) {
	return s.Impl.Quit(request.GroupId)
}

// ScanJoin 扫码进群
func (s Server) ScanJoin(ctx context.Context, request *ScanJoinRequest) (*JoinResult, error) {
	return s.Impl.ScanJoin(request.QrcodeUrl)
}

// ScanJoinEnterprise 企业微信扫码进群
func (s Server) ScanJoinEnterprise(ctx context.Context, request *ScanJoinRequest) (*JoinResult, error) {
	return s.Impl.ScanJoinEnterprise(request.QrcodeUrl)
}

// ConsentJoin 同意入群邀请
func (s Server) ConsentJoin(ctx context.Context, request *ConsentJoinRequest) (*JoinResult, error) {
	return s.Impl.ConsentJoin(request.InviteUrl)
}

// GetGroupByKey 按键查询缓存群信息
func (s Server) GetGroupByKey(ctx context.Context, request *GetGroupByKeyRequest) (*GetGroupResponse, error) {
	info, found := s.Impl.GetGroupByKey(request.Key)
	if !found {
		return &GetGroupResponse{Found: false}, nil
	}
	return &GetGroupResponse{Found: true, Group: info}, nil
}

// GetGroupByStrategy 按策略查询缓存群信息
func (s Server) GetGroupByStrategy(ctx context.Context, request *GetGroupByStrategyRequest) (*GetGroupResponse, error) {
	info, found := s.Impl.GetGroupByStrategy(request.Key, request.Strategy)
	if !found {
		return &GetGroupResponse{Found: false}, nil
	}
	return &GetGroupResponse{Found: true, Group: info}, nil
}

// GetGroupMembers 获取缓存群成员列表
func (s Server) GetGroupMembers(ctx context.Context, request *GetGroupMembersByKeyRequest) (*GetGroupMembersResponse, error) {
	members, found := s.Impl.GetGroupMembers(request.GroupId)
	if !found {
		return &GetGroupMembersResponse{Found: false}, nil
	}
	return &GetGroupMembersResponse{Found: true, Members: members}, nil
}
