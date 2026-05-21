// Package labelability 提供标签能力的实现。
package labelability

import (
	"github.com/sbgayhub/golem/host/api"
	sdk "github.com/sbgayhub/golem/sdk/label"

	labelapi "github.com/sbgayhub/golem/host/api/label"
)

// ability 标签能力实现（直连型）
type ability struct {
	api labelapi.LabelService
}

func init() {
	sdk.Instance = &ability{api: labelapi.Get()}
}

// List 获取标签列表
func (a ability) List() ([]*sdk.LabelPair, error) {
	resp, err := a.api.List()
	if resp == nil || err != nil {
		return nil, err
	}
	labels := make([]*sdk.LabelPair, 0, len(resp.List))
	for _, l := range resp.List {
		var pair sdk.LabelPair
		if err := api.TransformProto(l, &pair); err != nil {
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
	if resp.List != nil {
		if err := api.TransformProto(resp.List, &pair); err != nil {
			return nil, err
		}
	}
	return &pair, nil
}

// Delete 删除标签
func (a ability) Delete(labelIds string) (*sdk.Delete_Response, error) {
	resp, err := a.api.Delete(labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.Delete_Response
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新标签名称
func (a ability) Update(labelId uint32, name string) (*sdk.Update_Response, error) {
	resp, err := a.api.Update(labelId, name)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.Update_Response
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyContactLabels 修改联系人标签
func (a ability) ModifyContactLabels(usernames []string, labelIds string) (*sdk.ModifyContact_Response, error) {
	resp, err := a.api.ModifyContactLabels(usernames, labelIds)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.ModifyContact_Response
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
