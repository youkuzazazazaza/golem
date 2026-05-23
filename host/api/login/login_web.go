//go:build !lib

// Package loginapi 提供登录服务的 web 实现（通过 HTTP 调用远程服务）。
package loginapi

import (
	"sync"

	"github.com/sbgayhub/golem/host/api"
	baseapi "github.com/sbgayhub/golem/host/api/base"
)

// web 登录服务 web 实现
type web struct{}

// Get 获取 LoginService 单例（web 模式）
var Get = sync.OnceValue(func() LoginService {
	return &web{}
})

// Login 执行扫码登录
func (w web) Login() (*QRCodeResult, error) {
	var resp struct {
		UUID        string `json:"uuid"`
		Data        string `json:"data"`
		CheckTime   uint32 `json:"check_time"`
		ExpiredTime uint32 `json:"expired_time"`
	}
	if err := api.GetHttp().Get("/api/login/login").DoJson(&resp); err != nil {
		return nil, err
	}
	return &QRCodeResult{
		Qrcode:      &baseapi.Buffer{Data: []byte(resp.Data)},
		Uuid:        &resp.UUID,
		CheckTime:   &resp.CheckTime,
		ExpiredTime: &resp.ExpiredTime,
	}, nil
}

// Init 首次登录后初始化
func (w web) Init() (*InitResponse, error) {
	if _, err := api.GetHttp().Get("/api/login/init").Do(); err != nil {
		return nil, err
	}
	return &InitResponse{Code: 0, Message: "ok"}, nil
}

// Refresh 刷新登录状态
func (w web) Refresh() (*OperateResponse, error) {
	var status string
	if err := api.GetHttp().Get("/api/login/status").DoJson(&status); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0, Message: status}, nil
}

// Wakeup 唤醒登录
func (w web) Wakeup() (*OperateResponse, error) {
	if _, err := api.GetHttp().Get("/api/login/awaken").Do(); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0, Message: "ok"}, nil
}

// Logout 登出
func (w web) Logout() (*OperateResponse, error) {
	if _, err := api.GetHttp().Get("/api/login/logout").Do(); err != nil {
		return nil, err
	}
	return &OperateResponse{Code: 0, Message: "ok"}, nil
}

// PasswordLogin 使用账号密码登录
func (w web) PasswordLogin(req *PasswordLoginRequest) (*PasswordLoginResult, error) {
	var resp PasswordLoginResult
	if err := api.GetHttp().Post("/api/login/password").Body(map[string]any{
		"username":    req.GetUsername(),
		"password":    req.GetPassword(),
		"device_type": req.GetDeviceType(),
		"device_name": req.GetDeviceName(),
		"device_id":   req.GetDeviceId(),
	}).DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
