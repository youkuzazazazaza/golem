//go:build lib

// Package contactapi 提供联系人服务的 lib 实现（直接调用底层实现）。
package contactapi

import (
	"sync"

	"golem/pkg/contact"

	"github.com/sbgayhub/golem/host/api/util"
)

// lib 联系人服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 ContactService 单例（lib 模式）
var Get = sync.OnceValue(func() ContactService {
	return &lib{}
})

// List 获取联系人列表（增量同步）
func (l lib) List(contactSequence, groupSequence int32) (*ListContactsResponse, error) {
	resp, err := contact.List(contactSequence, groupSequence)
	if resp == nil || err != nil {
		return nil, err
	}
	var result ListContactsResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAll 获取全部联系人列表（分页查询）
func (l lib) ListAll(contactSequence, groupSequence, offset, limit int32) ([]*ContactInfo, error) {
	resp, err := contact.ListAll(contactSequence, groupSequence, offset, limit)
	if err != nil {
		return nil, err
	}
	result := make([]*ContactInfo, len(resp))
	for i, c := range resp {
		result[i] = &ContactInfo{
			Username: c.GetUsername(),
			Nickname: c.GetNickname(),
			Remark:   c.GetRemark(),
		}
	}
	return result, nil
}

// Detail 获取联系人详细信息
func (l lib) Detail(usernames, groups []string) (*GetContactDetailResponse, error) {
	resp, err := contact.Detail(usernames, groups)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetContactDetailResponse
	for _, mc := range resp.GetContactList() {
		var apiContact ModifyContact
		if err := util.TransformProto(mc, &apiContact); err != nil {
			return nil, err
		}
		result.ContactList = append(result.ContactList, &apiContact)
	}
	return &result, nil
}

// SetRemark 设置联系人备注
func (l lib) SetRemark(username, remark string) (*OperateResponse, error) {
	resp, err := contact.SetRemark(username, remark)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Search 搜索联系人
func (l lib) Search(keyword string, fromScene, searchScene uint32) (*SearchContactResponse, error) {
	resp, err := contact.Search(keyword, fromScene, searchScene)
	if resp == nil || err != nil {
		return nil, err
	}
	return &SearchContactResponse{
		Username: resp.GetUsername().GetValue(),
		Nickname: resp.GetNickname().GetValue(),
	}, nil
}

// Verify 通过好友验证
func (l lib) Verify(v1, v2 string, scene int) (*VerifyUserResponse, error) {
	resp, err := contact.Verify(v1, v2, scene)
	if resp == nil || err != nil {
		return nil, err
	}
	var result VerifyUserResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Request 发送好友申请
func (l lib) Request(v1, v2, content string, operate, scene int) (*VerifyUserResponse, error) {
	resp, err := contact.Request(v1, v2, content, operate, scene)
	if resp == nil || err != nil {
		return nil, err
	}
	var result VerifyUserResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BlacklistAdd 添加到黑名单
func (l lib) BlacklistAdd(username string) (*OperateResponse, error) {
	resp, err := contact.BlacklistAdd(username)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BlacklistRemove 从黑名单移除
func (l lib) BlacklistRemove(username string) (*OperateResponse, error) {
	resp, err := contact.BlacklistRemove(username)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete 删除联系人
func (l lib) Delete(username string) (*OperateResponse, error) {
	resp, err := contact.Delete(username)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LbsFind 附近的人
func (l lib) LbsFind(latitude, longitude float32, operate uint32) (*LbsFindResponse, error) {
	resp, err := contact.LbsFind(latitude, longitude, operate)
	if resp == nil || err != nil {
		return nil, err
	}
	var result LbsFindResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadContact 上传通讯录匹配好友
func (l lib) UploadContact(phones []string, currentPhone string, operate int32) (*UploadContactResponse, error) {
	resp, err := contact.UploadContact(phones, currentPhone, operate)
	if resp == nil || err != nil {
		return nil, err
	}
	var result UploadContactResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
