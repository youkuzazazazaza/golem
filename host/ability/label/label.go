// Package labelability 提供标签能力的实现。
package labelability

import (
	sdk "github.com/sbgayhub/golem/sdk/label"

	api "github.com/sbgayhub/golem/host/api/label"
	"github.com/sbgayhub/golem/host/api/util"
)

// ability 标签能力实现（直连型）
type ability struct {
	api api.LabelService
}

func init() {
	sdk.Instance = &ability{api: api.Get()}
}

// List 获取标签列表
func (a ability) List() ([]*sdk.LabelPair, error) {
	resp, err := a.api.List()
	if resp == nil || err != nil {
		return nil, err
	}
	labels := make([]*sdk.LabelPair, 0, len(resp.Labels))
	for _, l := range resp.Labels {
		var pair sdk.LabelPair
		if err := util.TransformProto(l, &pair); err != nil {
			return nil, err
		}
		labels = append(labels, &pair)
	}
	return labels, nil
}

// Add 添加标签
func (a ability) Add(name string) (*sdk.LabelPair, error) {
	resp, err := a.api.Add(name)
	if resp == nil || err != nil {
		return nil, err
	}
	var pair sdk.LabelPair
	if err := util.TransformProto(resp.Label, &pair); err != nil {
		return nil, err
	}
	return &pair, nil
}

// Delete 删除标签
func (a ability) Delete(labelIds string) (*sdk.OperateResponse, error) {
	resp, err := a.api.Delete(labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新标签名称
func (a ability) Update(labelId uint32, name string) (*sdk.OperateResponse, error) {
	resp, err := a.api.Update(labelId, name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyContactLabels 修改联系人标签
func (a ability) ModifyContactLabels(usernames []string, labelIds string) (*sdk.OperateResponse, error) {
	resp, err := a.api.ModifyContactLabels(usernames, labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.OperateResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
