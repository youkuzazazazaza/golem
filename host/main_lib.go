//go:build lib

package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golem"
	"golem/pkg/login"

	"github.com/phsym/console-slog"
	"github.com/sbgayhub/golem/host/ability"
	"github.com/sbgayhub/golem/host/ability/sync"
	"github.com/sbgayhub/golem/host/plugin"
)

func init() {
	handler := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelDebug,
		AddSource:  true,
		TimeFormat: "2006-01-02 15:04:05",
	})

	slog.SetDefault(slog.New(handler))
}

// lib模式
func main() {
	// 初始化协议层
	if err := golem.Initial(golem.WithSyncCallback(sync.CallBack), golem.WithLogger(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{})))); err != nil {
		slog.Warn("初始化协议层出错", "err", err)
		return
	}

	// 初始化能力层
	if err := ability.Initial(); err != nil {
		slog.Error("能力层注册失败", "err", err)
		return
	}

	// 初始化插件管理器
	if err := plugin.Initial(); err != nil {
		slog.Error("插件管理器初始化失败", "err", err)
		return
	}

	// 检查登录状态
	if user, err := login.Check(); err != nil {
		slog.Error("登录失败", "err", err)
		return
	} else {
		_ = user
		//contactability.SetUser(contact.User(user))
	}

	// 等待中断信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("正在关闭...")
	//ability.Destroy()
	plugin.Destroy()
	golem.Stop()
}
