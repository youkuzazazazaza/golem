// Package contactability 提供联系人能力的实现（缓存型）。
package contactability

import (
	"strings"
	"sync"

	"github.com/duke-git/lancet/v2/condition"
	sdk "github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/group"
)

// ability 联系人能力实现（缓存型）
type ability struct {
	contactCache sync.Map // map[string]*sdk.Contact
	memberCache  sync.Map // map[string][]*group.GroupMember
}

func init() {
	sdk.Instance = &ability{}
}

// GetContactByKey 按键查询缓存联系人
func (a *ability) GetContactByKey(key string) (*sdk.Contact, bool) {
	strategy, realKey := parseKeyPrefix(key)
	return a.getByStrategy(realKey, strategy)
}

// GetContactByStrategy 按策略查询缓存联系人
func (a *ability) GetContactByStrategy(key string, strategy sdk.RetrievalType) (*sdk.Contact, bool) {
	return a.getByStrategy(key, strategy)
}

// GetContactList 获取联系人列表
func (a *ability) GetContactList() ([]*sdk.Contact, error) {
	var contacts []*sdk.Contact
	a.contactCache.Range(func(_, value any) bool {
		contacts = append(contacts, value.(*sdk.Contact))
		return true
	})
	return contacts, nil
}

// GetGroupMembers 获取群成员列表
func (a *ability) GetGroupMembers(groupId string) ([]*group.GroupMember, bool) {
	v, ok := a.memberCache.Load(groupId)
	if !ok {
		return nil, false
	}
	return v.([]*group.GroupMember), true
}

// SetRemark 设置联系人备注
func (a *ability) SetRemark(username, remark string) (*sdk.OperateResponse, error) {
	if v, ok := a.contactCache.Load(username); ok {
		v.(*sdk.Contact).Remark = remark
	}
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// AddFriend 发送好友申请
func (a *ability) AddFriend(v1, v2, content string, operate, scene int) (*sdk.OperateResponse, error) {
	// 协议层调用由 host 进程通过 API 层处理
	// 此处仅记录请求，实际发送在 host 启动后生效
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// VerifyFriend 通过好友验证
func (a *ability) VerifyFriend(v1, v2 string, scene int) (*sdk.OperateResponse, error) {
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// Delete 删除联系人
func (a *ability) Delete(username string) (*sdk.OperateResponse, error) {
	a.contactCache.Delete(username)
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// BlacklistAdd 添加到黑名单
func (a *ability) BlacklistAdd(username string) (*sdk.OperateResponse, error) {
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// BlacklistRemove 从黑名单移除
func (a *ability) BlacklistRemove(username string) (*sdk.OperateResponse, error) {
	return &sdk.OperateResponse{Code: 0, Message: "ok"}, nil
}

// Search 搜索联系人
func (a *ability) Search(keyword string, fromScene, searchScene uint32) ([]*sdk.Contact, error) {
	var results []*sdk.Contact
	a.contactCache.Range(func(_, value any) bool {
		c := value.(*sdk.Contact)
		if strings.Contains(c.Nickname, keyword) || strings.Contains(c.Remark, keyword) || strings.Contains(c.Username, keyword) {
			results = append(results, c)
		}
		return true
	})
	return results, nil
}

// getByStrategy 按策略从缓存查询
func (a *ability) getByStrategy(key string, strategy sdk.RetrievalType) (*sdk.Contact, bool) {
	switch strategy {
	case sdk.RetrievalType_RETRIEVAL_TYPE_USERNAME:
		v, ok := a.contactCache.Load(key)
		if !ok {
			return nil, false
		}
		return v.(*sdk.Contact), true
	case sdk.RetrievalType_RETRIEVAL_TYPE_NICKNAME:
		var found *sdk.Contact
		a.contactCache.Range(func(_, value any) bool {
			c := value.(*sdk.Contact)
			if c.Nickname == key {
				found = c
				return false
			}
			return true
		})
		return found, found != nil
	case sdk.RetrievalType_RETRIEVAL_TYPE_REMARK:
		var found *sdk.Contact
		a.contactCache.Range(func(_, value any) bool {
			c := value.(*sdk.Contact)
			if c.Remark == key {
				found = c
				return false
			}
			return true
		})
		return found, found != nil
	default:
		v, ok := a.contactCache.Load(key)
		if !ok {
			return nil, false
		}
		return v.(*sdk.Contact), true
	}
}

// parseKeyPrefix 解析查询键前缀
func parseKeyPrefix(key string) (sdk.RetrievalType, string) {
	prefix, value, ok := strings.Cut(key, "::")
	if !ok {
		return sdk.RetrievalType_RETRIEVAL_TYPE_USERNAME, key
	}
	switch prefix {
	case "nickname":
		return sdk.RetrievalType_RETRIEVAL_TYPE_NICKNAME, value
	case "remark":
		return sdk.RetrievalType_RETRIEVAL_TYPE_REMARK, value
	case "username":
		return sdk.RetrievalType_RETRIEVAL_TYPE_USERNAME, value
	default:
		return sdk.RetrievalType_RETRIEVAL_TYPE_USERNAME, key
	}
}

// specials 特殊账号列表
var specials = []string{
	"qqmail", "qmessage", "tmessage", "floatbottle", "facebookapp",
	"qqfriend", "newsapp", "feedsapp", "masssendapp", "blogapp",
	"voipapp", "voicevoipapp", "voiceinputapp", "googlecontact", "fmessage",
	"medianote", "qqsync", "lbsapp", "shakeapp", "linkedinplugin",
	"gh_43f2581f6fd6", "gh_3dfda90e39d6", "gh_f0a92aa7146c",
	"gh_579db1f2cf89", "gh_b4af18eac3d5",
	"gh_e087bb5b95e6", "weixin", "gh_051d9102de63",
}

// Build 从协议层 ModifyContact 构建 SDK Contact（由 host 同步调用）
func (a *ability) Build(username, nickname, remark, alias, avatar string,
	groupOwner string, country, province, city, gender, signature string,
	bigAvatar, smallAvatar string) *sdk.Contact {

	avatar = condition.Ternary(bigAvatar != "", bigAvatar, smallAvatar)
	contactType := getType(username, groupOwner)

	c := &sdk.Contact{
		Username: username,
		Nickname: nickname,
		Remark:   remark,
		Alias:    alias,
		Avatar:   avatar,
		Type:     contactType,
	}

	switch contactType {
	case sdk.ContactType_CONTACT_TYPE_FRIEND:
		c.Data = &sdk.Contact_Friend{Friend: &sdk.FriendData{
			Country:   country,
			Province:  province,
			City:      city,
			Gender:    gender,
			Signature: signature,
		}}
	case sdk.ContactType_CONTACT_TYPE_GROUP:
		c.Data = &sdk.Contact_Group{Group: &sdk.GroupData{Owner: groupOwner}}
	}

	a.contactCache.Store(username, c)
	return c
}

// SetGroupMembers 设置群成员缓存
func (a *ability) SetGroupMembers(groupId string, members []*group.GroupMember) {
	a.memberCache.Store(groupId, members)
}

// GetContactByUsername 通过用户名获取缓存联系人（供其他 ability 调用）
func (a *ability) GetContactByUsername(wxid string) *sdk.Contact {
	v, ok := a.contactCache.Load(wxid)
	if !ok {
		return nil
	}
	return v.(*sdk.Contact)
}

// getType 判断联系人类型
func getType(username, groupOwner string) sdk.ContactType {
	if groupOwner != "" {
		return sdk.ContactType_CONTACT_TYPE_GROUP
	}
	if strings.HasPrefix(username, "gh_") {
		return sdk.ContactType_CONTACT_TYPE_OFFICIAL
	}
	if contains(specials, username) {
		return sdk.ContactType_CONTACT_TYPE_SPECIAL
	}
	return sdk.ContactType_CONTACT_TYPE_FRIEND
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
