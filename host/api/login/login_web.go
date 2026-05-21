//go:build !lib

// Package loginapi 提供登录服务的 web 实现（通过 HTTP 调用远程服务）。
package loginapi

import (
	"sync"
)

// web 登录服务 web 实现
type web struct{}

// Get 获取 LoginService 单例（web 模式）
var Get = sync.OnceValue(func() LoginService {
	return &web{}
})

// Login 执行扫码登录
func (w web) Login() (*QRCodeResult, error) {
	var resp QRCodeResult
	if err := api.GetHttp().Get("/api/login/login").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Init 首次登录后初始化
func (w web) Init() (*InitResponse, error) {
	var resp InitResponse
	if err := api.GetHttp().Get("/api/login/init").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Refresh 刷新登录状态
func (w web) Refresh() (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Get("/api/login/status").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Wakeup 唤醒登录
func (w web) Wakeup() (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Get("/api/login/awaken").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Logout 登出
func (w web) Logout() (*OperateResponse, error) {
	var resp OperateResponse
	if err := api.GetHttp().Get("/api/login/logout").DoProto(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
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
