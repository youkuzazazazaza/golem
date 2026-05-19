// Package util 提供 HTTP 客户端和通用工具函数。
package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/sbgayhub/golem/host/config"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Response 统一 API 响应格式
type Response struct {
	Code    int64           `json:"code"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

// Client HTTP 客户端
type Client struct {
	BaseURL string
	Token string
	HTTPCli *http.Client
}

// GetHttp 获取 HTTP 客户端单例
// 从配置中获取 URL，初始化 HTTP 客户端
var GetHttp = sync.OnceValue(func() *Client {
	cfg := config.Get()
	return &Client{
		BaseURL: cfg.URL,
		Token:   cfg.Token,
		HTTPCli: http.DefaultClient,
	}
})

// Get 发送 GET 请求
func (c *Client) Get(path string) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil, nil)
}

// Post 发送 POST 请求（JSON Body）
func (c *Client) Post(path string, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}
	return c.doRequest(http.MethodPost, path, nil, bytes.NewReader(jsonBody))
}

// Put 发送 PUT 请求（JSON Body）
func (c *Client) Put(path string, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}
	return c.doRequest(http.MethodPut, path, nil, bytes.NewReader(jsonBody))
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(path string) ([]byte, error) {
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(method, path string, header map[string]string, body io.Reader) ([]byte, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	// 解析统一响应格式
	var wrapper Response
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		// 如果不是统一格式，返回原始响应
		return respBody, nil
	}

	// 检查业务错误
	if wrapper.Code != 0 {
		return nil, errors.New(wrapper.Message)
	}

	return wrapper.Data, nil
}

// ParseResponse 解析响应数据到指定类型
func ParseResponse(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// ParseProtoResponse 解析响应数据到 proto 消息类型
func ParseProtoResponse(data []byte, m proto.Message) error {
	if len(data) == 0 {
		return nil
	}
	return protojson.Unmarshal(data, m)
}

// BuildQueryString 构建 URL 查询字符串
func BuildQueryString(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}
	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return "?" + strings.Join(parts, "&")
}
