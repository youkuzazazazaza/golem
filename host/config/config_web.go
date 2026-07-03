//go:build !lib

package config

type Config struct {
	Owner     string `toml:"owner" comment:"机器人所有者"`
	Forbidden string `toml:"forbidden" comment:"无权限提示"`
	URL       string `toml:"url" comment:"golem api 地址"`
	Token     string `toml:"token" comment:"golem api token"`
}
