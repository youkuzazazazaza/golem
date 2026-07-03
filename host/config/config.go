//go:build !lib

package config

import "sync"

var Get = sync.OnceValue(func() Config {
	return Config{
		URL:   "http://localhost:8085/api",
		Token: "",
	}
})
