//go:build !lib

// Package groupapi 提供群聊服务的 web 实现（通过 HTTP 调用远程服务）。
package groupapi

import (
	"fmt"
	"sync"

	"github.com/sbgayhub/golem/host/api/util"
)

// web 群聊服务 web 实现（通过 HTTP 调用远程服务）
type web struct{}

// Get 获取 GroupService 单例（web 模式）
var Get = sync.OnceValue(func() GroupService {
	return &web{}
})

// Create 创建群聊
func (w web) Create(members []string) (*CreateGroupResponse, error) {
	req := CreateGroupRequest{Members: members}
	data, err := util.GetHttp().Post("/groups", &req)
	if err != nil {
		return nil, err
	}
	var resp CreateGroupResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FacingCreate 面对面建群
func (w web) FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateGroupResponse, error) {
	req := FacingCreateGroupRequest{
		Password:  password,
		Latitude:  latitude,
		Longitude: longitude,
		Operate:   operate,
	}
	data, err := util.GetHttp().Post("/groups/facing", &req)
	if err != nil {
		return nil, err
	}
	var resp FacingCreateGroupResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInfo 获取群详细信息
func (w web) GetInfo(groupID string) (*GetGroupInfoDetailResponse, error) {
	data, err := util.GetHttp().Get(fmt.Sprintf("/groups/info/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp GetGroupInfoDetailResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMemberDetail 获取群成员详情
func (w web) GetMemberDetail(groupID string) (*GetGroupMemberDetailResponse, error) {
	data, err := util.GetHttp().Get(fmt.Sprintf("/groups/members/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp GetGroupMemberDetailResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetQRCode 获取群二维码
func (w web) GetQRCode(groupID string) (*GetGroupQRCodeResponse, error) {
	data, err := util.GetHttp().Get(fmt.Sprintf("/groups/qrcode/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp GetGroupQRCodeResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddMember 添加群成员
func (w web) AddMember(groupID string, members []string) (*AddGroupMemberResponse, error) {
	req := AddGroupMemberRequest{
		GroupId: groupID,
		Members: members,
	}
	data, err := util.GetHttp().Post(fmt.Sprintf("/groups/members/%s", groupID), &req)
	if err != nil {
		return nil, err
	}
	var resp AddGroupMemberResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// InviteMember 邀请群成员
func (w web) InviteMember(groupID string, members []string) (*InviteGroupMemberResponse, error) {
	req := InviteGroupMemberRequest{
		GroupId: groupID,
		Members: members,
	}
	data, err := util.GetHttp().Post(fmt.Sprintf("/groups/invite/%s", groupID), &req)
	if err != nil {
		return nil, err
	}
	var resp InviteGroupMemberResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveMember 移除群成员
func (w web) RemoveMember(groupID string, members []string) (*RemoveGroupMemberResponse, error) {
	data, err := util.GetHttp().Delete(fmt.Sprintf("/groups/members/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp RemoveGroupMemberResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetName 设置群名称
func (w web) SetName(groupID, name string) (*OperateResponse, error) {
	data, err := util.GetHttp().Put(fmt.Sprintf("/groups/name/%s?name=%s", groupID, name), nil)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetAnnouncement 设置群公告
func (w web) SetAnnouncement(groupID, content string) (*SetAnnouncementResponse, error) {
	req := SetAnnouncementRequest{
		GroupId: groupID,
		Content: content,
	}
	data, err := util.GetHttp().Put(fmt.Sprintf("/groups/announcement/%s", groupID), &req)
	if err != nil {
		return nil, err
	}
	var resp SetAnnouncementResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetRemark 设置群备注
func (w web) SetRemark(groupID, remark string) (*OperateResponse, error) {
	data, err := util.GetHttp().Put(fmt.Sprintf("/groups/remark/%s?remark=%s", groupID, remark), nil)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetContactList 保存到通讯录
func (w web) SetContactList(groupID string, save bool) (*OperateResponse, error) {
	saveStr := "false"
	if save {
		saveStr = "true"
	}
	data, err := util.GetHttp().Put(fmt.Sprintf("/groups/contact-list/%s?save=%s", groupID, saveStr), nil)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetAdmin 设置群管理员
func (w web) SetAdmin(groupID string, members []string) (*OperateResponse, error) {
	req := SetGroupAdminRequest{
		GroupId: groupID,
		Members: members,
	}
	data, err := util.GetHttp().Post(fmt.Sprintf("/groups/admins/%s", groupID), &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveAdmin 移除群管理员
func (w web) RemoveAdmin(groupID string, members []string) (*OperateResponse, error) {
	data, err := util.GetHttp().Delete(fmt.Sprintf("/groups/admins/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransferOwner 转让群主
func (w web) TransferOwner(groupID, newOwner string) (*OperateResponse, error) {
	data, err := util.GetHttp().Post(fmt.Sprintf("/groups/transfer/%s?new_owner=%s", groupID, newOwner), nil)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Quit 退出群聊
func (w web) Quit(groupID string) (*OperateResponse, error) {
	data, err := util.GetHttp().Delete(fmt.Sprintf("/groups/quit/%s", groupID))
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ScanJoin 扫码进群
func (w web) ScanJoin(qrcodeURL string) (*JoinResult, error) {
	req := ScanJoinRequest{QrcodeUrl: qrcodeURL}
	data, err := util.GetHttp().Post("/groups/join/scan", &req)
	if err != nil {
		return nil, err
	}
	var resp JoinResult
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ScanJoinEnterprise 企业微信扫码进群
func (w web) ScanJoinEnterprise(qrcodeURL string) (*JoinResult, error) {
	return w.ScanJoin(qrcodeURL)
}

// ConsentJoin 同意入群邀请
func (w web) ConsentJoin(inviteURL string) (*JoinResult, error) {
	req := ConsentJoinRequest{InviteUrl: inviteURL}
	data, err := util.GetHttp().Post("/groups/join/consent", &req)
	if err != nil {
		return nil, err
	}
	var resp JoinResult
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
