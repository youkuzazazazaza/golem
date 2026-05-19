package favor

import (
	"context"
)

// Client 实现 Ability 接口，通过 gRPC 调用远程收藏服务
type Client struct {
	Client FavorServiceClient
}

var _ Ability = (*Client)(nil)

// GetInfo 获取收藏容量信息
func (c Client) GetInfo() (*FavorInfo, error) {
	resp, err := c.Client.GetInfo(context.Background(), &GetInfoRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Info, nil
}

// GetItem 获取收藏项详情
func (c Client) GetItem(favId int32) ([]*FavorItem, error) {
	resp, err := c.Client.GetItem(context.Background(), &GetFavItemRequest{FavId: favId})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// BatchGetItems 批量获取收藏项
func (c Client) BatchGetItems(favIds []int32) ([]*FavorItem, error) {
	resp, err := c.Client.BatchGetItems(context.Background(), &BatchGetFavItemsRequest{FavIds: favIds})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// Delete 删除收藏项
func (c Client) Delete(favId int32) ([]*DeleteResult, error) {
	resp, err := c.Client.Delete(context.Background(), &DeleteFavItemRequest{FavId: favId})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// BatchDelete 批量删除收藏项
func (c Client) BatchDelete(favIds []int32) ([]*DeleteResult, error) {
	resp, err := c.Client.BatchDelete(context.Background(), &BatchDeleteFavItemsRequest{FavIds: favIds})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// Sync 同步收藏列表
func (c Client) Sync(key []byte) (*SyncFavorResponse, error) {
	return c.Client.Sync(context.Background(), &SyncFavorRequest{Key: key})
}

// Server 实现 FavorServiceServer 接口，将 gRPC 请求委托给 Ability 实现
type Server struct {
	UnimplementedFavorServiceServer
	Impl Ability
}

// GetInfo 获取收藏容量信息
func (s Server) GetInfo(ctx context.Context, request *GetInfoRequest) (*GetInfoResponse, error) {
	info, err := s.Impl.GetInfo()
	if err != nil {
		return nil, err
	}
	return &GetInfoResponse{Info: info}, nil
}

// GetItem 获取收藏项详情
func (s Server) GetItem(ctx context.Context, request *GetFavItemRequest) (*GetFavItemResponse, error) {
	items, err := s.Impl.GetItem(request.FavId)
	if err != nil {
		return nil, err
	}
	return &GetFavItemResponse{Items: items}, nil
}

// BatchGetItems 批量获取收藏项
func (s Server) BatchGetItems(ctx context.Context, request *BatchGetFavItemsRequest) (*BatchGetFavItemsResponse, error) {
	items, err := s.Impl.BatchGetItems(request.FavIds)
	if err != nil {
		return nil, err
	}
	return &BatchGetFavItemsResponse{Items: items}, nil
}

// Delete 删除收藏项
func (s Server) Delete(ctx context.Context, request *DeleteFavItemRequest) (*DeleteFavItemResponse, error) {
	results, err := s.Impl.Delete(request.FavId)
	if err != nil {
		return nil, err
	}
	return &DeleteFavItemResponse{Results: results}, nil
}

// BatchDelete 批量删除收藏项
func (s Server) BatchDelete(ctx context.Context, request *BatchDeleteFavItemsRequest) (*BatchDeleteFavItemsResponse, error) {
	results, err := s.Impl.BatchDelete(request.FavIds)
	if err != nil {
		return nil, err
	}
	return &BatchDeleteFavItemsResponse{Results: results}, nil
}

// Sync 同步收藏列表
func (s Server) Sync(ctx context.Context, request *SyncFavorRequest) (*SyncFavorResponse, error) {
	return s.Impl.Sync(request.Key)
}
