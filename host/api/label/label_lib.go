//go:build lib

// Package labelapi 提供标签服务的 lib 实现（直接调用底层实现）。
package labelapi

import (
	"sync"

	"golem/pkg/label"

	"github.com/sbgayhub/golem/host/api/util"
)

// lib 标签服务 lib 实现
type lib struct{}

// Get 获取 LabelService 单例（lib 模式）
var Get = sync.OnceValue(func() LabelService {
	return &lib{}
})

// List 获取标签列表
func (l lib) List() (*ListLabelsResponse, error) {
	resp, err := label.List()
	if resp == nil || err != nil {
		return nil, err
	}
	var result ListLabelsResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Add 添加标签
func (l lib) Add(name string) (*AddLabelResponse, error) {
	resp, err := label.Add(name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result AddLabelResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete 删除标签
func (l lib) Delete(labelIds string) (*OperateResponse, error) {
	resp, err := label.Delete(labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新标签名称
func (l lib) Update(labelId uint32, name string) (*OperateResponse, error) {
	resp, err := label.Update(labelId, name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyContactLabels 修改联系人标签
func (l lib) ModifyContactLabels(usernames []string, labelIds string) (*OperateResponse, error) {
	resp, err := label.ModifyContactLabels(usernames, labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
