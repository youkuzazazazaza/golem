//go:build lib

// Package favorapi 提供收藏服务的 lib 实现（直接调用底层实现）。
package favorapi

import (
	"sync"

	"golem/pkg/favor"

	"github.com/sbgayhub/golem/host/api/util"
)

// lib 收藏服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 FavorService 单例（lib 模式）
var Get = sync.OnceValue(func() FavorService {
	return &lib{}
})

// GetInfo 获取收藏容量信息
func (l lib) GetInfo() (*GetInfoResponse, error) {
	resp, err := favor.GetInfo()
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetInfoResponse
	if err := util.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetItem 获取收藏项详情
func (l lib) GetItem(favId int32) (*GetFavItemResponse, error) {
	resp, err := favor.GetItem(favId)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*FavorItem, 0, len(resp.GetObjects()))
	for _, obj := range resp.GetObjects() {
		var item FavorItem
		if err := util.TransformProto(obj, &item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return &GetFavItemResponse{Items: items}, nil
}

// BatchGetItems 批量获取收藏项
func (l lib) BatchGetItems(favIds []int32) (*BatchGetFavItemsResponse, error) {
	resp, err := favor.BatchGetItems(favIds)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*FavorItem, 0, len(resp.GetObjects()))
	for _, obj := range resp.GetObjects() {
		var item FavorItem
		if err := util.TransformProto(obj, &item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return &BatchGetFavItemsResponse{Items: items}, nil
}

// Delete 删除收藏项
func (l lib) Delete(favId int32) (*DeleteFavItemResponse, error) {
	resp, err := favor.Delete(favId)
	if resp == nil || err != nil {
		return nil, err
	}
	results := make([]*DeleteResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		var result DeleteResult
		if err := util.TransformProto(r, &result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return &DeleteFavItemResponse{Results: results}, nil
}

// BatchDelete 批量删除收藏项
func (l lib) BatchDelete(favIds []int32) (*BatchDeleteFavItemsResponse, error) {
	resp, err := favor.BatchDelete(favIds)
	if resp == nil || err != nil {
		return nil, err
	}
	results := make([]*DeleteResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		var result DeleteResult
		if err := util.TransformProto(r, &result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return &BatchDeleteFavItemsResponse{Results: results}, nil
}

// Sync 同步收藏列表
func (l lib) Sync(key []byte) (*SyncFavorResponse, error) {
	resp, err := favor.Sync(key)
	if resp == nil || err != nil {
		return nil, err
	}
	items := make([]*FavorItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		var favorItem FavorItem
		if err := util.TransformProto(item, &favorItem); err != nil {
			return nil, err
		}
		items = append(items, &favorItem)
	}
	return &SyncFavorResponse{
		Items:   items,
		Key:     resp.Key,
		HasMore: resp.HasMore,
	}, nil
}
