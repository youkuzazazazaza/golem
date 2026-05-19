package contact

import (
	"context"

	"github.com/sbgayhub/golem/sdk/group"
)

// Client 实现 Ability 接口，通过 gRPC 调用远程联系人服务
type Client struct {
	Client ContactServiceClient
}

var _ Ability = (*Client)(nil)

// GetContactByKey 按键查询缓存联系人
func (c Client) GetContactByKey(key string) (*Contact, bool) {
	resp, err := c.Client.GetContactByKey(context.Background(), &GetContactByKeyRequest{Key: key})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Contact, true
}

// GetContactByStrategy 按策略查询缓存联系人
func (c Client) GetContactByStrategy(key string, strategy RetrievalType) (*Contact, bool) {
	resp, err := c.Client.GetContactByStrategy(context.Background(), &GetContactByStrategyRequest{
		Key:      key,
		Strategy: strategy,
	})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Contact, true
}

// GetContactList 获取联系人列表
func (c Client) GetContactList() ([]*Contact, error) {
	resp, err := c.Client.GetContactList(context.Background(), &GetContactListRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Contacts, nil
}

// GetGroupMembers 获取群成员列表
func (c Client) GetGroupMembers(groupId string) ([]*group.GroupMember, bool) {
	resp, err := c.Client.GetGroupMembers(context.Background(), &GetGroupMembersRequest{GroupId: groupId})
	if err != nil || !resp.Found {
		return nil, false
	}
	return resp.Members, true
}

// SetRemark 设置联系人备注
func (c Client) SetRemark(username, remark string) (*OperateResponse, error) {
	return c.Client.SetRemark(context.Background(), &SetRemarkRequest{
		Username: username,
		Remark:   remark,
	})
}

// AddFriend 发送好友申请
func (c Client) AddFriend(v1, v2, content string, operate, scene int) (*OperateResponse, error) {
	return c.Client.AddFriend(context.Background(), &AddFriendRequest{
		V1:      v1,
		V2:      v2,
		Content: content,
		Operate: int32(operate),
		Scene:   int32(scene),
	})
}

// VerifyFriend 通过好友验证
func (c Client) VerifyFriend(v1, v2 string, scene int) (*OperateResponse, error) {
	return c.Client.VerifyFriend(context.Background(), &VerifyFriendRequest{
		V1:    v1,
		V2:    v2,
		Scene: int32(scene),
	})
}

// Delete 删除联系人
func (c Client) Delete(username string) (*OperateResponse, error) {
	return c.Client.Delete(context.Background(), &DeleteContactRequest{Username: username})
}

// BlacklistAdd 添加到黑名单
func (c Client) BlacklistAdd(username string) (*OperateResponse, error) {
	return c.Client.BlacklistAdd(context.Background(), &BlacklistRequest{Username: username})
}

// BlacklistRemove 从黑名单移除
func (c Client) BlacklistRemove(username string) (*OperateResponse, error) {
	return c.Client.BlacklistRemove(context.Background(), &BlacklistRequest{Username: username})
}

// Search 搜索联系人
func (c Client) Search(keyword string, fromScene, searchScene uint32) ([]*Contact, error) {
	resp, err := c.Client.Search(context.Background(), &SearchContactRequest{
		Keyword:     keyword,
		FromScene:   fromScene,
		SearchScene: searchScene,
	})
	if err != nil {
		return nil, err
	}
	return resp.Contacts, nil
}

// Server 实现 ContactServiceServer 接口，将 gRPC 请求委托给 Ability 实现
type Server struct {
	UnimplementedContactServiceServer
	Impl Ability
}

// GetContactByKey 按键查询缓存联系人
func (s Server) GetContactByKey(ctx context.Context, request *GetContactByKeyRequest) (*GetContactResponse, error) {
	contact, found := s.Impl.GetContactByKey(request.Key)
	return &GetContactResponse{Found: found, Contact: contact}, nil
}

// GetContactByStrategy 按策略查询缓存联系人
func (s Server) GetContactByStrategy(ctx context.Context, request *GetContactByStrategyRequest) (*GetContactResponse, error) {
	contact, found := s.Impl.GetContactByStrategy(request.Key, request.Strategy)
	return &GetContactResponse{Found: found, Contact: contact}, nil
}

// GetGroupMembers 获取群成员列表
func (s Server) GetGroupMembers(ctx context.Context, request *GetGroupMembersRequest) (*GetGroupMembersResponse, error) {
	members, found := s.Impl.GetGroupMembers(request.GroupId)
	return &GetGroupMembersResponse{Found: found, Members: members}, nil
}

// GetContactList 获取联系人列表
func (s Server) GetContactList(ctx context.Context, request *GetContactListRequest) (*GetContactListResponse, error) {
	contacts, err := s.Impl.GetContactList()
	if err != nil {
		return nil, err
	}
	return &GetContactListResponse{Contacts: contacts}, nil
}

// SetRemark 设置联系人备注
func (s Server) SetRemark(ctx context.Context, request *SetRemarkRequest) (*OperateResponse, error) {
	return s.Impl.SetRemark(request.Username, request.Remark)
}

// AddFriend 发送好友申请
func (s Server) AddFriend(ctx context.Context, request *AddFriendRequest) (*OperateResponse, error) {
	return s.Impl.AddFriend(request.V1, request.V2, request.Content, int(request.Operate), int(request.Scene))
}

// VerifyFriend 通过好友验证
func (s Server) VerifyFriend(ctx context.Context, request *VerifyFriendRequest) (*OperateResponse, error) {
	return s.Impl.VerifyFriend(request.V1, request.V2, int(request.Scene))
}

// Delete 删除联系人
func (s Server) Delete(ctx context.Context, request *DeleteContactRequest) (*OperateResponse, error) {
	return s.Impl.Delete(request.Username)
}

// BlacklistAdd 添加到黑名单
func (s Server) BlacklistAdd(ctx context.Context, request *BlacklistRequest) (*OperateResponse, error) {
	return s.Impl.BlacklistAdd(request.Username)
}

// BlacklistRemove 从黑名单移除
func (s Server) BlacklistRemove(ctx context.Context, request *BlacklistRequest) (*OperateResponse, error) {
	return s.Impl.BlacklistRemove(request.Username)
}

// Search 搜索联系人
func (s Server) Search(ctx context.Context, request *SearchContactRequest) (*SearchResponse, error) {
	contacts, err := s.Impl.Search(request.Keyword, request.FromScene, request.SearchScene)
	if err != nil {
		return nil, err
	}
	return &SearchResponse{Contacts: contacts}, nil
}
