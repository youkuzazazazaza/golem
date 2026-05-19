// Package favorability 提供收藏能力的实现。
package favorability

import (
	sdk "github.com/sbgayhub/golem/sdk/favor"

	api "github.com/sbgayhub/golem/host/api/favor"
	"github.com/sbgayhub/golem/host/api/util"
)

// ability 收藏能力实现（直连型）
type ability struct {
	api api.FavorService
}

func init() {
	sdk.Instance = &ability{api: api.Get()}
}

// GetInfo 获取收藏容量信息
func (a ability) GetInfo() (*sdk.FavorInfo, error) {
	resp, err := a.api.GetInfo()
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.FavorInfo
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetItem 获取收藏项详情
func (a ability) GetItem(favId int32) ([]*sdk.FavorItem, error) {
	resp, err := a.api.GetItem(favId)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*sdk.FavorItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		var sdkItem sdk.FavorItem
		if err := util.TransformProto(item, &sdkItem); err != nil {
			return nil, err
		}
		items = append(items, &sdkItem)
	}
	return items, nil
}

// BatchGetItems 批量获取收藏项
func (a ability) BatchGetItems(favIds []int32) ([]*sdk.FavorItem, error) {
	resp, err := a.api.BatchGetItems(favIds)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*sdk.FavorItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		var sdkItem sdk.FavorItem
		if err := util.TransformProto(item, &sdkItem); err != nil {
			return nil, err
		}
		items = append(items, &sdkItem)
	}
	return items, nil
}

// Delete 删除收藏项
func (a ability) Delete(favId int32) ([]*sdk.DeleteResult, error) {
	resp, err := a.api.Delete(favId)
	if resp == nil || err != nil {
		return nil, err
	}
	results := make([]*sdk.DeleteResult, 0, len(resp.Results))
	for _, r := range resp.Results {
		var sdkResult sdk.DeleteResult
		if err := util.TransformProto(r, &sdkResult); err != nil {
			return nil, err
		}
		results = append(results, &sdkResult)
	}
	return results, nil
}

// BatchDelete 批量删除收藏项
func (a ability) BatchDelete(favIds []int32) ([]*sdk.DeleteResult, error) {
	resp, err := a.api.BatchDelete(favIds)
	if resp == nil || err != nil {
		return nil, err
	}
	results := make([]*sdk.DeleteResult, 0, len(resp.Results))
	for _, r := range resp.Results {
		var sdkResult sdk.DeleteResult
		if err := util.TransformProto(r, &sdkResult); err != nil {
			return nil, err
		}
		results = append(results, &sdkResult)
	}
	return results, nil
}

// Sync 同步收藏列表
func (a ability) Sync(key []byte) (*sdk.SyncFavorResponse, error) {
	resp, err := a.api.Sync(key)
	if resp == nil || err != nil {
		return nil, err
	}
	var result sdk.SyncFavorResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
