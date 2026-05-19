//go:build !lib

// Package contactapi 提供联系人服务的 web 实现（通过 HTTP 调用远程服务）。
package contactapi

import (
	"sync"

	"github.com/sbgayhub/golem/host/api/util"
)

// web 联系人服务 web 实现（通过 HTTP 调用远程服务）
type web struct{}

// Get 获取 ContactService 单例（web 模式）
var Get = sync.OnceValue(func() ContactService {
	return &web{}
})

// List 获取联系人列表（增量同步）
func (w web) List(contactSequence, groupSequence int32) (*ListContactsResponse, error) {
	req := ListContactsRequest{
		ContactSequence: contactSequence,
		GroupSequence:   groupSequence,
	}
	data, err := util.GetHttp().Post("/contact/list", &req)
	if err != nil {
		return nil, err
	}
	var resp ListContactsResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAll 获取全部联系人列表（分页查询）
func (w web) ListAll(contactSequence, groupSequence, offset, limit int32) ([]*ContactInfo, error) {
	req := ListAllContactsRequest{
		ContactSequence: contactSequence,
		GroupSequence:   groupSequence,
		Offset:          offset,
		Limit:           limit,
	}
	data, err := util.GetHttp().Post("/contact/list_all", &req)
	if err != nil {
		return nil, err
	}
	var resp ListAllContactsResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return resp.Contacts, nil
}

// Detail 获取联系人详细信息
func (w web) Detail(usernames, groups []string) (*GetContactDetailResponse, error) {
	req := GetContactDetailRequest{
		Usernames: usernames,
		Groups:    groups,
	}
	data, err := util.GetHttp().Post("/contact/detail", &req)
	if err != nil {
		return nil, err
	}
	var resp GetContactDetailResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetRemark 设置联系人备注
func (w web) SetRemark(username, remark string) (*OperateResponse, error) {
	req := SetRemarkRequest{
		Username: username,
		Remark:   remark,
	}
	data, err := util.GetHttp().Post("/contact/set_remark", &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search 搜索联系人
func (w web) Search(keyword string, fromScene, searchScene uint32) (*SearchContactResponse, error) {
	req := SearchContactRequest{
		Keyword:     keyword,
		FromScene:   fromScene,
		SearchScene: searchScene,
	}
	data, err := util.GetHttp().Post("/contact/search", &req)
	if err != nil {
		return nil, err
	}
	var resp SearchContactResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify 通过好友验证
func (w web) Verify(v1, v2 string, scene int) (*VerifyUserResponse, error) {
	req := VerifyUserRequest{
		V1:    v1,
		V2:    v2,
		Scene: int32(scene),
	}
	data, err := util.GetHttp().Post("/contact/verify", &req)
	if err != nil {
		return nil, err
	}
	var resp VerifyUserResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Request 发送好友申请
func (w web) Request(v1, v2, content string, operate, scene int) (*VerifyUserResponse, error) {
	req := SendFriendRequest{
		V1:      v1,
		V2:      v2,
		Content: content,
		Operate: int32(operate),
		Scene:   int32(scene),
	}
	data, err := util.GetHttp().Post("/contact/request", &req)
	if err != nil {
		return nil, err
	}
	var resp VerifyUserResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlacklistAdd 添加到黑名单
func (w web) BlacklistAdd(username string) (*OperateResponse, error) {
	req := BlacklistRequest{
		Username: username,
	}
	data, err := util.GetHttp().Post("/contact/blacklist/add", &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlacklistRemove 从黑名单移除
func (w web) BlacklistRemove(username string) (*OperateResponse, error) {
	req := BlacklistRequest{
		Username: username,
	}
	data, err := util.GetHttp().Post("/contact/blacklist/remove", &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete 删除联系人
func (w web) Delete(username string) (*OperateResponse, error) {
	req := DeleteContactRequest{
		Username: username,
	}
	data, err := util.GetHttp().Post("/contact/delete", &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LbsFind 附近的人
func (w web) LbsFind(latitude, longitude float32, operate uint32) (*LbsFindResponse, error) {
	req := LbsFindRequest{
		Latitude:  latitude,
		Longitude: longitude,
		Operate:   operate,
	}
	data, err := util.GetHttp().Post("/contact/lbs_find", &req)
	if err != nil {
		return nil, err
	}
	var resp LbsFindResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadContact 上传通讯录匹配好友
func (w web) UploadContact(phones []string, currentPhone string, operate int32) (*UploadContactResponse, error) {
	req := UploadContactRequest{
		Phones:       phones,
		CurrentPhone: currentPhone,
		Operate:      operate,
	}
	data, err := util.GetHttp().Post("/contact/upload_contact", &req)
	if err != nil {
		return nil, err
	}
	var resp UploadContactResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
