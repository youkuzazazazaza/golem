//go:build !lib

package main

import (
	"log/slog"
	"os"

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

}
