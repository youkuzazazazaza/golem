//go:build !lib

// Package labelapi 提供标签服务的 web 实现（通过 HTTP 调用远程服务）。
package labelapi

import (
	"fmt"
	"sync"

	"github.com/sbgayhub/golem/host/api/util"
)

// web 标签服务 web 实现
type web struct{}

// Get 获取 LabelService 单例（web 模式）
var Get = sync.OnceValue(func() LabelService {
	return &web{}
})

// List 获取标签列表
func (w web) List() (*ListLabelsResponse, error) {
	data, err := util.GetHttp().Get("/labels")
	if err != nil {
		return nil, err
	}
	var resp ListLabelsResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add 添加标签
func (w web) Add(name string) (*AddLabelResponse, error) {
	req := AddLabelRequest{Name: name}
	data, err := util.GetHttp().Post("/labels", &req)
	if err != nil {
		return nil, err
	}
	var resp AddLabelResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete 删除标签
func (w web) Delete(labelIds string) (*OperateResponse, error) {
	data, err := util.GetHttp().Delete(fmt.Sprintf("/labels/%s", labelIds))
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update 更新标签名称
func (w web) Update(labelId uint32, name string) (*OperateResponse, error) {
	req := UpdateLabelRequest{LabelId: labelId, Name: name}
	data, err := util.GetHttp().Put(fmt.Sprintf("/labels/%d", labelId), &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyContactLabels 修改联系人标签
func (w web) ModifyContactLabels(usernames []string, labelIds string) (*OperateResponse, error) {
	req := ModifyContactLabelsRequest{Usernames: usernames, LabelIds: labelIds}
	data, err := util.GetHttp().Put(fmt.Sprintf("/contacts/labels/%s", usernames[0]), &req)
	if err != nil {
		return nil, err
	}
	var resp OperateResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
