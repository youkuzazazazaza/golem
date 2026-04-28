//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

// Proto 生成Protobuf文件
func Proto() error {
	fmt.Println("生成Protobuf文件...")
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			return sh.Run("protoc", "--go_out=.", "--go-grpc_out=.", "--go_opt=paths=source_relative", "--go-grpc_opt=paths=source_relative", path)
		}
		return nil
	})
}
