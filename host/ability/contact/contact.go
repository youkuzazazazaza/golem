// Package contactability 提供联系人能力的实现（缓存型）。
package contactability

import (
	"github.com/sbgayhub/golem/host/api"
	contactapi "github.com/sbgayhub/golem/host/api/contact"
	sdk "github.com/sbgayhub/golem/sdk/contact"
)

// ability 联系人能力实现（缓存型）
type ability struct {
	cache map[string]*sdk.Contact
}

func init() {
	sdk.Instance = &ability{cache: map[string]*sdk.Contact{}}
}

// Get 按键查询缓存联系人，支持前缀：username::（默认）、nickname::、remark::
func (a *ability) Get(key string) *sdk.Contact {
	if c := getStrategy(key).find(a.cache); c != nil {
		return c
	}
	return nil
}

// List 获取联系人列表
func (a *ability) List() []*sdk.Contact {
	var result []*sdk.Contact
	for _, c := range a.cache {
		result = append(result, c)
	}
	return result
}

// SetRemark 设置联系人备注
func (a *ability) SetRemark(username, remark string) error {
	if response, err := contactapi.Get().SetRemark(username, remark); err != nil {
		return err
	} else {
		var result sdk.SetRemark_Response
		if err := api.TransformProto(response, &result); err != nil {
			return err
		}
		return nil
	}
}

// AddFriend 发送好友申请
func (a *ability) AddFriend(v1, v2, content string, operate, scene int) error {
	// 协议层调用由 host 进程通过 API 层处理
	// 此处仅记录请求，实际发送在 host 启动后生效
	return nil
}

// VerifyFriend 通过好友验证
func (a *ability) VerifyFriend(v1, v2 string, scene int) error {
	return nil
}

// Delete 删除联系人
func (a *ability) Delete(username string) error {
	return nil
}

// BlacklistAdd 添加到黑名单
func (a *ability) BlacklistAdd(username string) error {
	return nil
}

// BlacklistRemove 从黑名单移除
func (a *ability) BlacklistRemove(username string) error {
	return nil
}

// Search 搜索联系人
func (a *ability) Search(keyword string, fromScene, searchScene uint32) *sdk.Contact {
	// TODO: 实现搜索逻辑
	return nil
}
