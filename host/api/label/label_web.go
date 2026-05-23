//go:build !lib

// Package labelapi 提供标签服务的 web 实现（通过 HTTP 调用远程服务）。
package labelapi

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sbgayhub/golem/host/api"
)

// web 标签服务 web 实现
type web struct{}

// Get 获取 LabelService 单例（web 模式）
var Get = sync.OnceValue(func() LabelService {
	return &web{}
})

// List 获取标签列表
func (w web) List() (*ListLabelsResponse, error) {
	var resp ListLabelsResponse
	if err := api.GetHttp().Get("/api/labels").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add 添加标签
func (w web) Add(name string) (*AddLabelResponse, error) {
	var resp AddLabelResponse
	if err := api.GetHttp().Post("/api/labels").Body(map[string]any{"label_name": name}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete 删除标签
func (w web) Delete(labelIds string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Delete(fmt.Sprintf("/api/labels/%s", labelIds)).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update 更新标签名称
func (w web) Update(labelId uint32, name string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Put(fmt.Sprintf("/api/labels/%d", labelId)).Body(map[string]any{
		"label_name": name,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyContactLabels 修改联系人标签
func (w web) ModifyContactLabels(usernames []string, labelIds string) (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Put(fmt.Sprintf("/api/contacts/%s/labels", strings.Join(usernames, ","))).Body(map[string]any{
		"label_ids": labelIds,
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
