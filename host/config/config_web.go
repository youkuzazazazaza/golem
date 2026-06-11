//go:build !lib

package config

type Config struct {
	Owner string `toml:"owner" comment:"机器人所有者"`
	URL   string
	Token string
}
