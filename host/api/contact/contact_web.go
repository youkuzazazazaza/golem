//go:build !lib

// Package contactapi 提供联系人服务的 web 实现（通过 HTTP 调用远程服务）。
package contactapi

import (
	"fmt"
	"sync"

	"github.com/sbgayhub/golem/host/api"
)

// web 联系人服务 web 实现（通过 HTTP 调用远程服务）
type web struct{}

// Get 获取 ContactService 单例（web 模式）
var Get = sync.OnceValue(func() ContactService {
	return &web{}
})

// List 获取联系人列表（增量同步）
func (w web) List() ([]string, error) {
	var resp []string
	if err := api.GetHttp().Get("/api/contacts").DoJson(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Detail 获取联系人详细信息
func (w web) Detail(usernames []string) (*GetContactDetailResponse, error) {
	var resp GetContactDetailResponse
	if err := api.GetHttp().Post("/api/contacts/detail").Body(usernames).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetRemark 设置联系人备注
func (w web) SetRemark(username, remark string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Put(fmt.Sprintf("/api/contacts/remark/%s", username)).Body(map[string]any{
		"remark": remark,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search 搜索联系人
func (w web) Search(keyword string, fromScene, searchScene uint32) (*SearchContactResponse, error) {
	var resp SearchContactResponse
	if err := api.GetHttp().Post("/api/contacts/search").Body(map[string]any{
		"keyword":      keyword,
		"from_scene":   fromScene,
		"search_scene": searchScene,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify 通过好友验证
func (w web) Verify(v1, v2 string, scene int) (*VerifyUserResponse, error) {
	var resp VerifyUserResponse
	if err := api.GetHttp().Post("/api/contacts/friend-requests/verify").Body(map[string]any{
		"v1":    v1,
		"v2":    v2,
		"scene": int32(scene),
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Request 发送好友申请
func (w web) Request(v1, v2, content string, operate, scene int) (*VerifyUserResponse, error) {
	var resp VerifyUserResponse
	if err := api.GetHttp().Post("/api/contacts/friend-requests").Body(map[string]any{
		"v1":      v1,
		"v2":      v2,
		"content": content,
		"operate": int32(operate),
		"scene":   int32(scene),
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlacklistAdd 添加到黑名单
func (w web) BlacklistAdd(username string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Post(fmt.Sprintf("/api/contacts/blacklist/%s", username)).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlacklistRemove 从黑名单移除
func (w web) BlacklistRemove(username string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Delete(fmt.Sprintf("/api/contacts/blacklist/%s", username)).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete 删除联系人
func (w web) Delete(username string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Delete(fmt.Sprintf("/api/contacts/%s", username)).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LbsFind 附近的人
func (w web) LbsFind(latitude, longitude float32, operate uint32) (*LbsResponse, error) {
	var resp LbsResponse
	if err := api.GetHttp().Post("/api/contacts/lbs").Body(map[string]any{
		"latitude":  latitude,
		"longitude": longitude,
		"operate":   operate,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadContact 上传通讯录匹配好友
func (w web) UploadContact(phones []string, currentPhone string, operate int32) (*UploadContactResponse, error) {
	var resp UploadContactResponse
	if err := api.GetHttp().Post("/api/contacts/upload").Body(map[string]any{
		"phones":        phones,
		"current_phone": currentPhone,
		"operate":       operate,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
