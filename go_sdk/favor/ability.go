package favor

// Ability 收藏能力接口（供插件嵌入使用）
type Ability interface {
	// GetInfo 获取收藏容量信息
	GetInfo() (*FavorInfo, error)
	// GetItem 获取收藏项详情
	GetItem(favId int32) ([]*FavorItem, error)
	// BatchGetItems 批量获取收藏项
	BatchGetItems(favIds []int32) ([]*FavorItem, error)
	// Delete 删除收藏项
	Delete(favId int32) ([]*DeleteResult, error)
	// BatchDelete 批量删除收藏项
	BatchDelete(favIds []int32) ([]*DeleteResult, error)
	// Sync 同步收藏列表（key 为空时从头同步，否则增量同步）
	Sync(key []byte) (*SyncFavorResponse, error)
}

// Instance 收藏能力实例（由 host/ability 层注入）
var Instance Ability
