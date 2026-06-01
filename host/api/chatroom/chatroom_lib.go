//go:build lib

// Package chatroomapi 提供群聊服务的 lib 实现（直接调用底层实现）。
package chatroomapi

import (
	"sync"

	"golem/pkg/chatroom"

	"github.com/sbgayhub/golem/host/api"
)

// lib 群聊服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 ChatroomService 单例（lib 模式）
var Get = sync.OnceValue(func() ChatroomService {
	return &lib{}
})

// Create 创建群聊
func (l lib) Create(members []string) (*CreateChatroomResponse, error) {
	resp, err := chatroom.Create(members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result CreateChatroomResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FacingCreate 面对面建群
func (l lib) FacingCreate(password string, latitude, longitude float32, operate uint32) (*FacingCreateChatroomResponse, error) {
	resp, err := chatroom.FacingCreate(password, latitude, longitude, operate)
	if resp == nil || err != nil {
		return nil, err
	}
	var result FacingCreateChatroomResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInfo 获取群详细信息
func (l lib) GetInfo(chatroomID string) (*GetChatroomInfoDetailResponse, error) {
	resp, err := chatroom.GetInfo(chatroomID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetChatroomInfoDetailResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (l lib) ListMembers(chatroomID string) (*ListMembersResponse, error) {
	resp, err := chatroom.ListMembers(chatroomID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result ListMembersResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (l lib) GetMemberDetail(chatroomID string, members []string) (*GetChatroomMembersResponse, error) {
	detail, err := chatroom.GetMemberDetail(chatroomID, members)
	if detail == nil || err != nil {
		return nil, err
	}
	var result GetChatroomMembersResponse
	if err := api.TransformProto(detail, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQRCode 获取群二维码
func (l lib) GetQRCode(chatroomID string) (*GetChatroomQRCodeResponse, error) {
	resp, err := chatroom.GetQRCode(chatroomID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetChatroomQRCodeResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddMember 添加群成员
func (l lib) AddMember(chatroomID string, members []string) (*AddChatroomMemberResponse, error) {
	resp, err := chatroom.AddMember(chatroomID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result AddChatroomMemberResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InviteMember 邀请群成员
func (l lib) InviteMember(chatroomID string, members []string) (*InviteChatroomMemberResponse, error) {
	resp, err := chatroom.InviteMember(chatroomID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result InviteChatroomMemberResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveMember 移除群成员
func (l lib) RemoveMember(chatroomID string, members []string) (*RemoveChatroomMemberResponse, error) {
	resp, err := chatroom.RemoveMember(chatroomID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result RemoveChatroomMemberResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetName 设置群名称
func (l lib) SetName(chatroomID, name string) (*OperateResponse, error) {
	resp, err := chatroom.SetName(chatroomID, name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAnnouncement 设置群公告
func (l lib) SetAnnouncement(chatroomID, content string) (*SetAnnouncementResponse, error) {
	resp, err := chatroom.SetAnnouncement(chatroomID, content)
	if resp == nil || err != nil {
		return nil, err
	}
	var result SetAnnouncementResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetRemark 设置群备注
func (l lib) SetRemark(chatroomID, remark string) (*OperateResponse, error) {
	resp, err := chatroom.SetRemark(chatroomID, remark)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetContactList 保存到通讯录
func (l lib) SetContactList(chatroomID string, save bool) (*OperateResponse, error) {
	resp, err := chatroom.SetContactList(chatroomID, save)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAdmin 设置群管理员
func (l lib) SetAdmin(chatroomID string, members []string) (*ChatroomAdminResponse, error) {
	resp, err := chatroom.SetAdmin(chatroomID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result ChatroomAdminResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveAdmin 移除群管理员
func (l lib) RemoveAdmin(chatroomID string, members []string) (*ChatroomAdminResponse, error) {
	resp, err := chatroom.RemoveAdmin(chatroomID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result ChatroomAdminResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TransferOwner 转让群主
func (l lib) TransferOwner(chatroomID, newOwner string) (*ChatroomAdminResponse, error) {
	resp, err := chatroom.TransferOwner(chatroomID, newOwner)
	if resp == nil || err != nil {
		return nil, err
	}
	var result ChatroomAdminResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Quit 退出群聊
func (l lib) Quit(chatroomID string) (*OperateResponse, error) {
	resp, err := chatroom.Quit(chatroomID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanJoin 扫码进群
func (l lib) ScanJoin(qrcodeURL string) (*JoinResult, error) {
	resp, err := chatroom.ScanJoin(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		ChatroomId: resp.ChatroomID,
		Message:    resp.Message,
	}, nil
}

// ScanJoinEnterprise 企业微信扫码进群
func (l lib) ScanJoinEnterprise(qrcodeURL string) (*JoinResult, error) {
	resp, err := chatroom.ScanJoinEnterprise(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		ChatroomId: resp.ChatroomID,
		Message:    resp.Message,
	}, nil
}

// ConsentJoin 同意入群邀请
func (l lib) ConsentJoin(inviteURL string) (*JoinResult, error) {
	resp, err := chatroom.ConsentJoin(inviteURL)
	if resp == nil || err != nil {
		return nil, err
	}
	return &JoinResult{
		ChatroomId: resp.ChatroomID,
		Message:    resp.Message,
	}, nil
}
