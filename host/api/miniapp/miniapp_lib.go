//go:build lib

// Package miniappapi 提供小程序服务的 lib 实现（直接调用底层实现）。
package miniappapi

import (
	"sync"

	"golem/pkg/miniapp"

	"github.com/sbgayhub/golem/host/api"
)

// lib 小程序服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 MiniAppService 单例（lib 模式）
var Get = sync.OnceValue(func() MiniAppService {
	return &lib{}
})

// JSLogin 小程序 JS 登录授权
func (l lib) JSLogin(appID string) (*JSLoginResponse, error) {
	resp, err := miniapp.JSLogin(appID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result JSLoginResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// QrcodeAuthLogin 小程序扫码授权登录
func (l lib) QrcodeAuthLogin(uuid string) (*StringResponse, error) {
	code, err := miniapp.QrcodeAuthLogin(uuid)
	if err != nil {
		return nil, err
	}
	return &StringResponse{Value: code}, nil
}

// GetRuntimeSession 获取小程序运行时会话
func (l lib) GetRuntimeSession(appID string) (*GetRuntimeSessionResponse, error) {
	resp, err := miniapp.GetRuntimeSession(appID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetRuntimeSessionResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSessionQRCode 获取小程序会话二维码
func (l lib) GetSessionQRCode(appID string) (*StringResponse, error) {
	url, err := miniapp.GetSessionQRCodeSimple(appID)
	if err != nil {
		return nil, err
	}
	return &StringResponse{Value: url}, nil
}

// OperateWxData 操作小程序数据
func (l lib) OperateWxData(appID string, data []byte, opt int32) (*JSOperateResponse, error) {
	resp, err := miniapp.OperateWxData(appID, data, opt)
	if resp == nil || err != nil {
		return nil, err
	}
	var result JSOperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloudCallFunction 调用小程序云函数
func (l lib) CloudCallFunction(appID string, data []byte) (*JSOperateResponse, error) {
	resp, err := miniapp.CloudCallFunction(appID, data)
	if resp == nil || err != nil {
		return nil, err
	}
	var result JSOperateResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendVerifyCode 发送手机验证码
func (l lib) SendVerifyCode(appID, mobile string, opcode int) (*PostVerifyCodeResponse, error) {
	resp, err := miniapp.SendVerifyCode(appID, mobile, opcode)
	if resp == nil || err != nil {
		return nil, err
	}
	var result PostVerifyCodeResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckVerifyCode 校验手机验证码
func (l lib) CheckVerifyCode(appID, mobile, code string, opcode int) (*CheckVerifyCodeResponse, error) {
	resp, err := miniapp.CheckVerifyCode(appID, mobile, code, opcode)
	if resp == nil || err != nil {
		return nil, err
	}
	var result CheckVerifyCodeResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddMobile 绑定手机号
func (l lib) AddMobile(appID, mobile, code string) (*CheckVerifyCodeResponse, error) {
	resp, err := miniapp.AddMobile(appID, mobile, code)
	if resp == nil || err != nil {
		return nil, err
	}
	var result CheckVerifyCodeResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DelMobile 解绑手机号
func (l lib) DelMobile(appID, mobile string) (*OperateResponse, error) {
	if err := miniapp.DelMobile(appID, mobile); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0}, nil
}

// GetAllMobile 获取已绑定手机号列表
func (l lib) GetAllMobile(appID string) (*OperateResponse, error) {
	if _, err := miniapp.GetAllMobile(appID); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0}, nil
}

// GetRandomAvatar 获取随机头像
func (l lib) GetRandomAvatar(appID string) (*OAuthGetRandomAvatarResponse, error) {
	resp, err := miniapp.GetRandomAvatar(appID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OAuthGetRandomAvatarResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddAvatar 设置头像
func (l lib) AddAvatar(appID, nickname, afilekey string) (*OperateResponse, error) {
	if err := miniapp.AddAvatar(appID, nickname, afilekey); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0}, nil
}

// UploadAvatarImg 上传自定义头像图片
func (l lib) UploadAvatarImg(appID, jpglink string) (*OAuthAddAvatarImageResponse, error) {
	resp, err := miniapp.UploadAvatarImg(appID, jpglink)
	if resp == nil || err != nil {
		return nil, err
	}
	var result OAuthAddAvatarImageResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRecord 获取小程序使用记录
func (l lib) GetRecord() (*GetUsageRecordResponse, error) {
	resp, err := miniapp.GetRecord()
	if resp == nil || err != nil {
		return nil, err
	}
	var result GetUsageRecordResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddRecord 添加小程序使用记录
func (l lib) AddRecord(username string) (*OperateResponse, error) {
	if err := miniapp.AddRecord(username); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0}, nil
}

// GetUserOpenID 获取用户 OpenID
func (l lib) GetUserOpenID(appID, username string) (*BizJsApiGetUserOpenIdResponse, error) {
	resp, err := miniapp.GetUserOpenID(appID, username)
	if resp == nil || err != nil {
		return nil, err
	}
	var result BizJsApiGetUserOpenIdResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OauthSdkApp SDK OAuth 确认授权
func (l lib) OauthSdkApp(appID string, scope []string, operate uint32) (*SDKOAuthAuthorizeConfirmResponse, error) {
	resp, err := miniapp.OauthSdkAppSimple(appID)
	if resp == nil || err != nil {
		return nil, err
	}
	var result SDKOAuthAuthorizeConfirmResponse
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ThirdAppGrant 第三方 APP OAuth 授权
func (l lib) ThirdAppGrant(appID, oauthURL string) (*StringResponse, error) {
	code, err := miniapp.ThirdAppGrantSimple(appID, oauthURL)
	if err != nil {
		return nil, err
	}
	return &StringResponse{Value: code}, nil
}
