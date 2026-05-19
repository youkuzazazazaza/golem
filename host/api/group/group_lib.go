//go:build lib

// Package groupapi 提供群聊服务的 lib 实现（直接调用底层实现）。
package groupapi

import (
	"sync"

	"golem/pkg/group"

	"github.com/sbgayhub/golem/host/api/util"
)

// lib 群聊服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 GroupService 单例（lib 模式）
var Get = sync.OnceValue(func() GroupService {
	return &lib{}
})

// Create 创建群聊
func (l lib) Create(members []string) (*CreateGroupResponse, error) {
	resp, err := group.Create(members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result CreateGroupResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FacingCreate 面对面建群
func (l lib) FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateGroupResponse, error) {
	resp, err := group.FacingCreate(password, latitude, longitude, operate)
	if resp == nil || err != nil {
		return nil, err
	}
	var result FacingCreateGroupResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInfo 获取群详细信息
func (l lib) GetInfo(groupID string) (*GetGroupInfoDetailResponse, error) {
	resp, err := group.GetInfo(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetGroupInfoDetailResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMemberDetail 获取群成员详情
func (l lib) GetMemberDetail(groupID string) (*GetGroupMemberDetailResponse, error) {
	resp, err := group.GetMemberDetail(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetGroupMemberDetailResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQRCode 获取群二维码
func (l lib) GetQRCode(groupID string) (*GetGroupQRCodeResponse, error) {
	resp, err := group.GetQRCode(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetGroupQRCodeResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddMember 添加群成员
func (l lib) AddMember(groupID string, members []string) (*AddGroupMemberResponse, error) {
	resp, err := group.AddMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result AddGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InviteMember 邀请群成员
func (l lib) InviteMember(groupID string, members []string) (*InviteGroupMemberResponse, error) {
	resp, err := group.InviteMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result InviteGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveMember 移除群成员
func (l lib) RemoveMember(groupID string, members []string) (*RemoveGroupMemberResponse, error) {
	resp, err := group.RemoveMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result RemoveGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetName 设置群名称
func (l lib) SetName(groupID, name string) (*OperateResponse, error) {
	resp, err := group.SetName(groupID, name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAnnouncement 设置群公告
func (l lib) SetAnnouncement(groupID, content string) (*SetAnnouncementResponse, error) {
	resp, err := group.SetAnnouncement(groupID, content)
	if resp == nil || err != nil {
		return nil, err
	}
	var result SetAnnouncementResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetRemark 设置群备注
func (l lib) SetRemark(groupID, remark string) (*OperateResponse, error) {
	resp, err := group.SetRemark(groupID, remark)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetContactList 保存到通讯录
func (l lib) SetContactList(groupID string, save bool) (*OperateResponse, error) {
	resp, err := group.SetContactList(groupID, save)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAdmin 设置群管理员
func (l lib) SetAdmin(groupID string, members []string) (*OperateResponse, error) {
	resp, err := group.SetAdmin(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveAdmin 移除群管理员
func (l lib) RemoveAdmin(groupID string, members []string) (*OperateResponse, error) {
	resp, err := group.RemoveAdmin(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TransferOwner 转让群主
func (l lib) TransferOwner(groupID, newOwner string) (*OperateResponse, error) {
	resp, err := group.TransferOwner(groupID, newOwner)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Quit 退出群聊
func (l lib) Quit(groupID string) (*OperateResponse, error) {
	resp, err := group.Quit(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanJoin 扫码进群
func (l lib) ScanJoin(qrcodeURL string) (*JoinResult, error) {
	resp, err := group.ScanJoin(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		GroupId: resp.GroupID,
		Message: resp.Message,
	}, nil
}

// ScanJoinEnterprise 企业微信扫码进群
func (l lib) ScanJoinEnterprise(qrcodeURL string) (*JoinResult, error) {
	resp, err := group.ScanJoinEnterprise(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		GroupId: resp.GroupID,
		Message: resp.Message,
	}, nil
}

// ConsentJoin 同意入群邀请
func (l lib) ConsentJoin(inviteURL string) (*JoinResult, error) {
	resp, err := group.ConsentJoin(inviteURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		GroupId: resp.GroupID,
		Message: resp.Message,
	}, nil
}
