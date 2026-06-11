//go:build !lib

package config

import "sync"

var Get = sync.OnceValue(func() HostConfig {
	return HostConfig{
		URL:   "http://localhost:8085/api",
		Token: "",
	}
})
