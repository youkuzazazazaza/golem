//go:build lib

package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"golem/config"

	"github.com/pelletier/go-toml/v2"
)

type HostConfig struct {
	Owner     string               `toml:"owner" comment:"机器人所有者"`
	Forbidden string               `toml:"forbidden" comment:"无权限提示"`
	Core      config.CoreConfig    `toml:"core" comment:"核心配置"`
	Server    config.ServerConfig  `toml:"server" comment:"Web服务配置"`
	Device    config.DeviceConfig  `toml:"device" comment:"设备配置"`
	Storage   config.StorageConfig `toml:"storage" comment:"存储配置（路径固定为./data，不可配置）"`
	Log       config.LogConfig     `toml:"log" comment:"日志配置"`
}

var path = "./data/config.toml"

// Get 获取配置文件
var Get = sync.OnceValue(func() *HostConfig {
	config, err := loadOrCreate()
	if err != nil {
		slog.Error("配置文件出现问题", "err", err)
		return nil
	}
	return config
})

// Save 保存配置到文件
func Save() error {
	return save(Get())
}

// 内部使用
func save(cfg *HostConfig) error {
	marshal, _ := json.Marshal(cfg)
	_ = marshal
	t := reflect.ValueOf(cfg).Elem().Type()
	slog.Info(fmt.Sprintf("[%s] %s", t.Name(), t.Field(0).Tag.Get("comment")))
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// defaults 返回默认配置
func defaults() *HostConfig {
	return &HostConfig{
		Core: config.CoreConfig{
			PrintBanner:         true,
			IgnoreHistory:       true,
			IgnoreTimeout:       10,
			RandomDelay:         true,
			RandomDelayMin:      500,
			RandomDelayMax:      3000,
			RandomDelayStrategy: "pending",
			QrcodeApi:           "https://api.qrtool.cn/?text=",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		Device: config.DeviceConfig{
			Type:       "ipad",
			KeyVersion: 144,
		},
		Log: config.LogConfig{
			Level:      "info",
			Output:     "both",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
			Compress:   true,
		},
	}
}

// load 从文件加载配置
func load() (*HostConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := defaults()
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// loadOrCreate 加载配置，如果文件不存在则创建默认配置并等待用户修改
func loadOrCreate() (*HostConfig, error) {
	// 检查./data目录是否存在
	if _, err := os.Stat(filepath.Dir(path)); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(path), 0644); err != nil {
			return nil, err
		}
	}

	cfg, err := load()
	if err == nil {
		return cfg, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	// 配置文件不存在，创建默认配置
	if err := save(defaults()); err != nil {
		return nil, fmt.Errorf("创建默认配置失败: %w", err)
	}

	fmt.Println("========================================")
	fmt.Println("配置文件不存在，已生成默认配置：", path)
	fmt.Println("请修改配置后重新启动服务")
	fmt.Println("========================================")
	fmt.Print("按回车键退出...")

	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)

	return nil, nil
}
