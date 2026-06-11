# Golem Plugins 插件合集

- example：示例插件

## 创建一个插件

插件是一个独立的 Go 可执行程序，通过 `github.com/sbgayhub/golem/sdk/plugin` 暴露给宿主。最小插件只需要实现 `plugin.Plugin` 接口；如果需要处理事件、命令或被其他插件调用，再按需实现对应的小接口。

### 1. 创建插件目录和模块

建议每个插件独立放在 `plugins/<plugin-name>` 目录下：

```bash
mkdir plugins/my_plugin
cd plugins/my_plugin
go mod init golem_plugin_my_plugin
```

如果插件目录在本仓库内开发，可以在 `go.mod` 中使用本地 SDK：

```go
module golem_plugin_my_plugin

go 1.25.0

require github.com/sbgayhub/golem/sdk v0.0.0

replace github.com/sbgayhub/golem/sdk => ../../go_sdk
```

### 2. 导入依赖

按需导入 SDK 包。`plugin` 是必需依赖；其他能力包只在插件实际使用时导入。

```go
import (
	"log/slog"
	"strings"

	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)
```

常用包说明：

- `plugin`：插件启动、元数据、事件、命令、配置、会话等基础能力。
- `message`：发送、转发、撤回、下载消息。
- `contact`：联系人查询、列表、备注、好友操作等。
- 其他能力包按需导入，例如 `chatroom`、`label`、`miniapp`、`cdn`、`favor`、`login`。

### 3. 定义配置和插件结构体

插件结构体用于声明插件自身状态，并通过字段嵌入或声明要使用的能力。宿主会按字段类型注入能力实例。

```go
type Config struct {
	Prefix string `toml:"prefix" comment:"回复前缀"`
}

type MyPlugin struct {
	plugin.ConfigAbility[Config]

	message message.Ability
	contact contact.Ability
	session plugin.SessionAbility
}
```

注意事项：

- `plugin.ConfigAbility[T]` 用于声明插件配置，并提供 `SaveConfig` 保存配置。
- `message.Ability`、`contact.Ability`、`plugin.SessionAbility` 等字段表示插件需要宿主注入对应能力。
- 字段可以命名声明；匿名嵌入也支持，但命名字段更直观，适合多个能力并存的场景。
- 只声明当前插件真实使用的能力，避免无意义依赖。

### 4. 实现基础插件接口

所有插件都必须实现 `plugin.Plugin` 接口，也就是 `GetMetadata` 方法。

```go
func (p *MyPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "my_plugin",
		Author:      "your-name",
		Version:     "1.0.0",
		Description: "my plugin",
		Next:        false,
		Priority:    0,
		AlwaysRun:   false,
	}
}
```

`Name` 会作为插件唯一标识，并对应 `plugins/config.toml` 中的配置段。

### 5. 按需实现事件插件接口

如果插件需要接收事件，实现 `plugin.EventPlugin`：

```go
func (p *MyPlugin) GetSubscriptions() []string {
	return []string{"message::text"}
}

func (p *MyPlugin) OnEvent(e *plugin.Event) (bool, error) {
	msg, ok := e.Payload.(*plugin.Event_Message)
	if !ok {
		return false, nil
	}

	content := strings.TrimSpace(msg.Message.Content)
	if content == "" {
		return false, nil
	}

	reply := p.Config.Prefix + content
	_, err := p.message.Send(&message.Message{
		Receiver: &contact.Contact{Username: e.Sender},
		Content:  reply,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
```

返回值约定：

- `bool` 表示事件是否已被当前插件处理。
- `error` 返回非空时表示处理失败，宿主会记录错误。

### 6. 按需实现命令插件接口

如果插件需要响应命令，实现 `plugin.CommandPlugin`。推荐使用结构体标签声明命令参数，再交给 SDK 绑定和分发。

```go
type EchoCommand struct {
	_      struct{} `cmd:"my echo" help:"回显文本" usage:"/my echo <text> [--upper]" example:"/my echo hello --upper"`
	Text   string   `arg:"text" help:"回显内容" required:"true" variadic:"true"`
	Upper  bool     `flag:"upper" help:"转换为大写" value:"false"`
}

func (p *MyPlugin) GetCommands() []string {
	return plugin.CommandCommands()
}

func (p *MyPlugin) OnCommand(cmd *plugin.Command) (string, error) {
	return plugin.DispatchCommand(cmd)
}

func (p *MyPlugin) handleEcho(cmd EchoCommand) (string, error) {
	text := cmd.Text
	if cmd.Upper {
		text = strings.ToUpper(text)
	}
	return text, nil
}
```

命令处理函数需要在 `main` 中注册：

```go
if err := plugin.RegisterCommand(p.handleEcho); err != nil {
	slog.Error("注册命令失败", "err", err)
	return
}
```

### 7. 按需实现被调用插件接口

如果插件需要给其他插件或宿主暴露能力，实现 `plugin.CalledPlugin`：

```go
import "fmt"

func (p *MyPlugin) GetCapabilities() []string {
	return []string{"echo"}
}

func (p *MyPlugin) OnCall(capability string, args map[string]string) (string, []byte, error) {
	switch capability {
	case "echo":
		return "text", []byte(args["text"]), nil
	default:
		return "", nil, fmt.Errorf("unsupported capability: %s", capability)
	}
}
```

### 8. 启动插件

在 `main` 函数中创建插件实例，初始化默认配置，注册命令，然后调用 `plugin.Start`。

```go
func main() {
	p := &MyPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: Config{
				Prefix: "echo: ",
			},
		},
	}

	if err := plugin.RegisterCommand(p.handleEcho); err != nil {
		slog.Error("注册命令失败", "err", err)
		return
	}

	plugin.Start(p)
}
```

### 9. 完整示例

```go
package main

import (
	"log/slog"
	"strings"

	"github.com/sbgayhub/golem/sdk/contact"
	"github.com/sbgayhub/golem/sdk/message"
	"github.com/sbgayhub/golem/sdk/plugin"
)

type Config struct {
	Prefix string `toml:"prefix" comment:"回复前缀"`
}

type MyPlugin struct {
	plugin.ConfigAbility[Config]

	message message.Ability
}

type EchoCommand struct {
	_     struct{} `cmd:"my echo" help:"回显文本" usage:"/my echo <text> [--upper]" example:"/my echo hello --upper"`
	Text  string   `arg:"text" help:"回显内容" required:"true" variadic:"true"`
	Upper bool     `flag:"upper" help:"转换为大写" value:"false"`
}

func (p *MyPlugin) GetMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Name:        "my_plugin",
		Author:      "your-name",
		Version:     "1.0.0",
		Description: "my plugin",
		Next:        false,
		Priority:    0,
		AlwaysRun:   false,
	}
}

func (p *MyPlugin) GetSubscriptions() []string {
	return []string{"message::text"}
}

func (p *MyPlugin) OnEvent(e *plugin.Event) (bool, error) {
	msg, ok := e.Payload.(*plugin.Event_Message)
	if !ok {
		return false, nil
	}

	content := strings.TrimSpace(msg.Message.Content)
	if content == "" {
		return false, nil
	}

	_, err := p.message.Send(&message.Message{
		Receiver: &contact.Contact{Username: e.Sender},
		Content:  p.Config.Prefix + content,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *MyPlugin) GetCommands() []string {
	return plugin.CommandCommands()
}

func (p *MyPlugin) OnCommand(cmd *plugin.Command) (string, error) {
	return plugin.DispatchCommand(cmd)
}

func (p *MyPlugin) handleEcho(cmd EchoCommand) (string, error) {
	text := cmd.Text
	if cmd.Upper {
		text = strings.ToUpper(text)
	}
	return p.Config.Prefix + text, nil
}

func main() {
	p := &MyPlugin{
		ConfigAbility: plugin.ConfigAbility[Config]{
			Config: Config{
				Prefix: "echo: ",
			},
		},
	}

	if err := plugin.RegisterCommand(p.handleEcho); err != nil {
		slog.Error("注册命令失败", "err", err)
		return
	}

	plugin.Start(p)
}
```

### 10. 配置插件

在 `plugins/config.toml` 中添加同名配置段：

```toml
[my_plugin]
enable = true
mode = 'blacklist'
limits = []

[my_plugin.config]
prefix = 'echo: '
```

### 11. 编译

在插件目录中执行：

```bash
go mod tidy
go build -o my_plugin.exe .
```

宿主加载插件时需要拿到编译后的可执行文件路径。插件名、配置段名和 `GetMetadata().Name` 建议保持一致，降低配置和排错成本。
