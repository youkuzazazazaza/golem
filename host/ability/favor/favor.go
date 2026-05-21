// Package favorability 提供收藏能力的实现。
package favorability

import (
	sdk "github.com/sbgayhub/golem/sdk/favor"

	"github.com/sbgayhub/golem/host/api"
	favorapi "github.com/sbgayhub/golem/host/api/favor"
)

// ability 收藏能力实现（直连型）
type ability struct {
	api favorapi.FavorService
}

func init() {
	sdk.Instance = &ability{api: favorapi.Get()}
}

// GetInfo 获取收藏容量信息
func (a ability) GetInfo() (*sdk.GetInfo_Info, error) {
	resp, err := a.api.GetInfo()
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.GetInfo_Info
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetItem 获取收藏项详情
func (a ability) GetItem(favId int32) ([]*sdk.Item, error) {
	resp, err := a.api.GetItem(favId)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*sdk.Item, 0, len(resp.Items))
	for _, item := range resp.Items {
		var sdkItem sdk.Item
		if err := api.TransformProto(item, &sdkItem); err != nil {
			return nil, err
		}
		items = append(items, &sdkItem)
	}
	return items, nil
}

// Delete 删除收藏项
func (a ability) Delete(favId int32) error {
	resp, err := a.api.Delete(favId)
	if resp == nil || err != nil {
		return err
	}
	results := make([]*sdk.Delete_Result, 0, len(resp.Results))
	for _, r := range resp.Results {
		var sdkResult sdk.Delete_Result
		if err := api.TransformProto(r, &sdkResult); err != nil {
			return err
		}
		results = append(results, &sdkResult)
	}
	return nil
}

// Sync 同步收藏列表
func (a ability) Sync(key []byte) (*sdk.Sync_Response, error) {
	resp, err := a.api.Sync(key)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.Sync_Response
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
