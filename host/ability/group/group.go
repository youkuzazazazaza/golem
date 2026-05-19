// Package groupability 提供群聊能力的实现。
package groupability

import (
	"strings"
	"sync"

	sdk "github.com/sbgayhub/golem/sdk/group"

	api "github.com/sbgayhub/golem/host/api/group"
	"github.com/sbgayhub/golem/host/api/util"
)

// ability 群聊能力实现（缓存型）
type ability struct {
	api         api.GroupService
	groupCache  sync.Map // map[string]*sdk.Group
	memberCache sync.Map // map[string][]*sdk.GroupMember
}

func init() {
	sdk.Instance = &ability{api: api.Get()}
}

// Create 创建群聊
func (a *ability) Create(members []string) (*sdk.CreateGroupResponse, error) {
	resp, err := a.api.Create(members)
	if resp == nil || err != nil {
		return nil, err
	}
	// 写入群缓存
	if resp.Group != "" {
		a.groupCache.Store(resp.Group, &sdk.Group{
			GroupId:        resp.Group,
			Name:           resp.Topic,
			BigAvatarUrl:   resp.BigAvatarUrl,
			SmallAvatarUrl: resp.SmallAvatarUrl,
		})
	}
	var result sdk.CreateGroupResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FacingCreate 面对面建群
func (a *ability) FacingCreate(password string, latitude, longitude float32, operate uint32) (*sdk.FacingCreateGroupResponse, error) {
	resp, err := a.api.FacingCreate(password, latitude, longitude, operate)
	if resp == nil || err != nil {
		return nil, err
	}
	if resp.Group != "" {
		a.groupCache.Store(resp.Group, &sdk.Group{
			GroupId: resp.Group,
		})
	}
	var result sdk.FacingCreateGroupResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInfo 获取群信息（缓存优先）
func (a *ability) GetInfo(groupID string) (*sdk.Group, error) {
	resp, err := a.api.GetInfo(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	info := &sdk.Group{
		GroupId:      groupID,
		Announcement: resp.Announcement,
	}
	a.groupCache.Store(groupID, info)
	return info, nil
}

// GetMemberDetail 获取群成员详情（缓存优先）
func (a *ability) GetMemberDetail(groupID string) ([]*sdk.GroupMember, error) {
	resp, err := a.api.GetMemberDetail(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	members := make([]*sdk.GroupMember, 0, len(resp.Members))
	for _, m := range resp.Members {
		var sdkMember sdk.GroupMember
		if err := util.TransformProto(m, &sdkMember); err != nil {
			return nil, err
		}
		members = append(members, &sdkMember)
	}
	a.memberCache.Store(groupID, members)
	return members, nil
}

// GetQRCode 获取群二维码
func (a *ability) GetQRCode(groupID string) (*sdk.GetGroupQRCodeResponse, error) {
	resp, err := a.api.GetQRCode(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.GetGroupQRCodeResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddMember 添加群成员
func (a *ability) AddMember(groupID string, members []string) (*sdk.AddGroupMemberResponse, error) {
	resp, err := a.api.AddMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.AddGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InviteMember 邀请群成员
func (a *ability) InviteMember(groupID string, members []string) (*sdk.InviteGroupMemberResponse, error) {
	resp, err := a.api.InviteMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.InviteGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveMember 移除群成员
func (a *ability) RemoveMember(groupID string, members []string) (*sdk.RemoveGroupMemberResponse, error) {
	resp, err := a.api.RemoveMember(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.RemoveGroupMemberResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetName 设置群名称（更新缓存）
func (a *ability) SetName(groupID, name string) (*sdk.OperateResponse, error) {
	resp, err := a.api.SetName(groupID, name)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		v.(*sdk.Group).Name = name
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAnnouncement 设置群公告（更新缓存）
func (a *ability) SetAnnouncement(groupID, content string) (*sdk.SetAnnouncementResponse, error) {
	resp, err := a.api.SetAnnouncement(groupID, content)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		v.(*sdk.Group).Announcement = content
	}
	var result sdk.SetAnnouncementResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetRemark 设置群备注（更新缓存）
func (a *ability) SetRemark(groupID, remark string) (*sdk.OperateResponse, error) {
	resp, err := a.api.SetRemark(groupID, remark)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		v.(*sdk.Group).Remark = remark
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetContactList 保存到通讯录
func (a *ability) SetContactList(groupID string, save bool) (*sdk.OperateResponse, error) {
	resp, err := a.api.SetContactList(groupID, save)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAdmin 设置群管理员（更新缓存）
func (a *ability) SetAdmin(groupID string, members []string) (*sdk.OperateResponse, error) {
	resp, err := a.api.SetAdmin(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		info := v.(*sdk.Group)
		existing := make(map[string]bool, len(info.Admins))
		for _, admin := range info.Admins {
			existing[admin] = true
		}
		for _, m := range members {
			if !existing[m] {
				info.Admins = append(info.Admins, m)
			}
		}
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveAdmin 移除群管理员（更新缓存）
func (a *ability) RemoveAdmin(groupID string, members []string) (*sdk.OperateResponse, error) {
	resp, err := a.api.RemoveAdmin(groupID, members)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		info := v.(*sdk.Group)
		removeSet := make(map[string]bool, len(members))
		for _, m := range members {
			removeSet[m] = true
		}
		filtered := info.Admins[:0]
		for _, admin := range info.Admins {
			if !removeSet[admin] {
				filtered = append(filtered, admin)
			}
		}
		info.Admins = filtered
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TransferOwner 转让群主（更新缓存）
func (a *ability) TransferOwner(groupID, newOwner string) (*sdk.OperateResponse, error) {
	resp, err := a.api.TransferOwner(groupID, newOwner)
	if resp == nil || err != nil {
		return nil, err
	}
	if v, ok := a.groupCache.Load(groupID); ok {
		v.(*sdk.Group).Owner = newOwner
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Quit 退出群聊（清除缓存）
func (a *ability) Quit(groupID string) (*sdk.OperateResponse, error) {
	resp, err := a.api.Quit(groupID)
	if resp == nil || err != nil {
		return nil, err
	}
	a.groupCache.Delete(groupID)
	a.memberCache.Delete(groupID)
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanJoin 扫码进群
func (a *ability) ScanJoin(qrcodeURL string) (*sdk.JoinResult, error) {
	resp, err := a.api.ScanJoin(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	if resp.GroupId != "" {
		a.groupCache.Store(resp.GroupId, &sdk.Group{GroupId: resp.GroupId})
	}
	return &sdk.JoinResult{
		GroupId: resp.GroupId,
		Message: resp.Message,
	}, nil
}

// ScanJoinEnterprise 企业微信扫码进群
func (a *ability) ScanJoinEnterprise(qrcodeURL string) (*sdk.JoinResult, error) {
	resp, err := a.api.ScanJoinEnterprise(qrcodeURL)
	if resp == nil || err != nil {
		return nil, err
	}
	if resp.GroupId != "" {
		a.groupCache.Store(resp.GroupId, &sdk.Group{GroupId: resp.GroupId})
	}
	return &sdk.JoinResult{
		GroupId: resp.GroupId,
		Message: resp.Message,
	}, nil
}

// ConsentJoin 同意入群邀请
func (a *ability) ConsentJoin(inviteURL string) (*sdk.JoinResult, error) {
	resp, err := a.api.ConsentJoin(inviteURL)
	if resp == nil || err != nil {
		return nil, err
	}
	if resp.GroupId != "" {
		a.groupCache.Store(resp.GroupId, &sdk.Group{GroupId: resp.GroupId})
	}
	return &sdk.JoinResult{
		GroupId: resp.GroupId,
		Message: resp.Message,
	}, nil
}

// GetGroupByKey 按键查询缓存群信息
func (a *ability) GetGroupByKey(key string) (*sdk.Group, bool) {
	strategy, realKey := parseKeyPrefix(key)
	return a.getByStrategy(realKey, strategy)
}

// GetGroupByStrategy 按策略查询缓存群信息
func (a *ability) GetGroupByStrategy(key string, strategy sdk.RetrievalType) (*sdk.Group, bool) {
	return a.getByStrategy(key, strategy)
}

// GetGroupMembers 获取缓存群成员列表
func (a *ability) GetGroupMembers(groupID string) ([]*sdk.GroupMember, bool) {
	v, ok := a.memberCache.Load(groupID)
	if !ok {
		return nil, false
	}
	return v.([]*sdk.GroupMember), true
}

// getByStrategy 按策略从缓存查询群信息
func (a *ability) getByStrategy(key string, strategy sdk.RetrievalType) (*sdk.Group, bool) {
	switch strategy {
	case sdk.RetrievalType_RETRIEVAL_TYPE_GROUP_ID:
		v, ok := a.groupCache.Load(key)
		if !ok {
			return nil, false
		}
		return v.(*sdk.Group), true
	case sdk.RetrievalType_RETRIEVAL_TYPE_NAME:
		var found *sdk.Group
		a.groupCache.Range(func(_, value any) bool {
			info := value.(*sdk.Group)
			if info.Name == key {
				found = info
				return false
			}
			return true
		})
		return found, found != nil
	default:
		v, ok := a.groupCache.Load(key)
		if !ok {
			return nil, false
		}
		return v.(*sdk.Group), true
	}
}

// parseKeyPrefix 解析查询键前缀
func parseKeyPrefix(key string) (sdk.RetrievalType, string) {
	prefix, value, ok := strings.Cut(key, "::")
	if !ok {
		return sdk.RetrievalType_RETRIEVAL_TYPE_GROUP_ID, key
	}
	switch prefix {
	case "name":
		return sdk.RetrievalType_RETRIEVAL_TYPE_NAME, value
	case "group_id":
		return sdk.RetrievalType_RETRIEVAL_TYPE_GROUP_ID, value
	default:
		return sdk.RetrievalType_RETRIEVAL_TYPE_GROUP_ID, key
	}
}
