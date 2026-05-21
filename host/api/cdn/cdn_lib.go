//go:build lib

// Package cdnapi 提供 CDN 文件上传/下载服务的 lib 实现（直接调用底层实现）。
package cdnapi

import (
	"bytes"
	"io"
	"sync"

	"golem/pkg/cdn"

	"github.com/sbgayhub/golem/host/api"
)

// lib CDN 服务 lib 实现（直接调用底层实现）
type lib struct{}

// Get 获取 CDNService 单例（lib 模式）
var Get = sync.OnceValue(func() CDNService {
	return &lib{}
})

// UploadImage CDN 上传聊天图片
func (l lib) UploadImage(receiver string, reader io.Reader) (*UploadImageResult, error) {
	resp, err := cdn.UploadImage(receiver, reader)
	if resp == nil || err != nil {
		return nil, err
	}
	var result UploadImageResult
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadMomentsImage CDN 上传朋友圈图片
func (l lib) UploadMomentsImage(imageData []byte) (*UploadSnsImageResult, error) {
	resp, err := cdn.UploadMomentsImage(bytes.NewReader(imageData))
	if resp == nil || err != nil {
		return nil, err
	}
	var result UploadSnsImageResult
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadVideo CDN 上传聊天视频
func (l lib) UploadVideo(receiver string, thumb []byte, reader io.Reader, duration uint32) (*UploadVideoResult, error) {
	resp, err := cdn.UploadVideo(receiver, reader, bytes.NewReader(thumb), duration)
	if resp == nil || err != nil {
		return nil, err
	}
	var result UploadVideoResult
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadMomentsVideo CDN 上传朋友圈视频
func (l lib) UploadMomentsVideo(videoData, thumbData []byte) (*UploadSnsVideoResult, error) {
	resp, err := cdn.UploadMomentsVideo(bytes.NewReader(videoData), bytes.NewReader(thumbData))
	if resp == nil || err != nil {
		return nil, err
	}
	var result UploadSnsVideoResult
	if err := api.TransformProto(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadImage CDN 下载高清图片（返回 ReadCloser 流）
func (l lib) DownloadImage(fileID, fileAesKey string) (io.ReadCloser, error) {
	data, err := cdn.DownloadImage(fileID, fileAesKey)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// DownloadVideoCover CDN 下载视频封面
func (l lib) DownloadVideoCover(fileID, fileAesKey string) ([]byte, error) {
	return cdn.DownloadVideoCover(fileID, fileAesKey)
}

// DownloadVideo CDN 下载聊天视频（返回 ReadCloser 流）
func (l lib) DownloadVideo(fileID, fileAesKey string) (io.ReadCloser, error) {
	data, err := cdn.DownloadVideo(fileID, fileAesKey)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// DownloadSnsVideo CDN 下载朋友圈视频
func (l lib) DownloadSnsVideo(videoURL string, encKey uint64) ([]byte, error) {
	return cdn.DownloadMomentsVideo(videoURL, encKey)
}
