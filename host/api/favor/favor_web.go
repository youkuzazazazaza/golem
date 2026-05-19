//go:build !lib

// Package favorapi 提供收藏服务的 web 实现（通过 HTTP 调用远程服务）。
package favorapi

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/sbgayhub/golem/host/api/util"
)

// web 收藏服务 web 实现（通过 HTTP 调用远程服务）
type web struct{}

// Get 获取 FavorService 单例（web 模式）
var Get = sync.OnceValue(func() FavorService {
	return &web{}
})

// GetInfo 获取收藏容量信息
func (w web) GetInfo() (*GetInfoResponse, error) {
	data, err := util.GetHttp().Get("/favor/info")
	if err != nil {
		return nil, err
	}
	var resp GetInfoResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetItem 获取收藏项详情
func (w web) GetItem(favId int32) (*GetFavItemResponse, error) {
	data, err := util.GetHttp().Get(fmt.Sprintf("/favor/item/%d", favId))
	if err != nil {
		return nil, err
	}
	var resp GetFavItemResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchGetItems 批量获取收藏项（web 模式逐个查询聚合）
func (w web) BatchGetItems(favIds []int32) (*BatchGetFavItemsResponse, error) {
	items := make([]*FavorItem, 0, len(favIds))
	for _, favId := range favIds {
		resp, err := w.GetItem(favId)
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
	}
	return &BatchGetFavItemsResponse{Items: items}, nil
}

// Delete 删除收藏项
func (w web) Delete(favId int32) (*DeleteFavItemResponse, error) {
	data, err := util.GetHttp().Delete(fmt.Sprintf("/favor/item/%d", favId))
	if err != nil {
		return nil, err
	}
	var resp DeleteFavItemResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchDelete 批量删除收藏项（web 模式逐个删除聚合）
func (w web) BatchDelete(favIds []int32) (*BatchDeleteFavItemsResponse, error) {
	results := make([]*DeleteResult, 0, len(favIds))
	for _, favId := range favIds {
		resp, err := w.Delete(favId)
		if err != nil {
			return nil, err
		}
		results = append(results, resp.Results...)
	}
	return &BatchDeleteFavItemsResponse{Results: results}, nil
}

// Sync 同步收藏列表
func (w web) Sync(key []byte) (*SyncFavorResponse, error) {
	path := "/favor/sync"
	if len(key) > 0 {
		path += "?key=" + base64.StdEncoding.EncodeToString(key)
	}
	data, err := util.GetHttp().Post(path, nil)
	if err != nil {
		return nil, err
	}
	var resp SyncFavorResponse
	if err := util.ParseProtoResponse(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
