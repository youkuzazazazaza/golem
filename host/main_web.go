//go:build !lib

package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	chatroomability "github.com/sbgayhub/golem/host/ability/chatroom"
	contactability "github.com/sbgayhub/golem/host/ability/contact"
	messageability "github.com/sbgayhub/golem/host/ability/message"
	hc "github.com/sbgayhub/golem/host/config"
	loginapi "github.com/sbgayhub/golem/host/api/login"
	userapi "github.com/sbgayhub/golem/host/api/user"
	"github.com/sbgayhub/golem/host/ability"
	"github.com/sbgayhub/golem/host/plugin"
	sdkcontact "github.com/sbgayhub/golem/sdk/contact"

	"github.com/phsym/console-slog"
)

func init() {
	handler := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelDebug,
		AddSource:  true,
		TimeFormat: "2006-01-02 15:04:05",
	})

	slog.SetDefault(slog.New(handler))
	slog.Info("日志系统初始化完成")
}

// web模式
func main() {
	// 加载配置
	cfg := hc.Get()
	slog.Info("web模式启动", "url", cfg.URL)
	messageability.SetOutboundReady(false)

	// 唤醒远程登录
	if _, err := loginapi.Get().Wakeup(); err != nil {
		slog.Error("唤醒登录失败", "err", err)
		return
	}

	// 获取用户信息
	profile, err := userapi.Get().GetProfile()
	if err != nil {
		slog.Error("获取用户信息失败", "err", err)
		return
	}
	info := profile.GetUserInfo()
	ext := profile.GetUserInfoExt()
	contactability.SetSelf(&sdkcontact.SelfInfo{
		Username: info.GetUserName().GetValue(),
		Nickname: info.GetNickName().GetValue(),
		Alias:    info.GetAlias(),
		Avatar:   ext.GetSmallAvatarUrl(),
		Uin:      info.GetUin(),
		Email:    info.GetEmail().GetValue(),
		Mobile:   info.GetMobile().GetValue(),
	})

	// 初始化联系人能力，加载联系人缓存
	contactability.Initial()
	// 初始化群组能力，加载群组缓存
	chatroomability.Initial()
	messageability.SetOutboundReady(true)

	// 初始化插件管理器
	if err := plugin.Initial(); err != nil {
		slog.Error("插件管理器初始化失败", "err", err)
		return
	}

	// 等待中断信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("正在关闭...")
	ability.Destroy()
	plugin.Destroy()
}
