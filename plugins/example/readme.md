# Example 插件 - 完整示例

这是一个**完整的示例插件**，展示了 Golem 插件系统的所有核心功能。如果你是第一次开发 Golem 插件，强烈建议从这个示例开始学习。

## 📚 目录

- [功能概览](#功能概览)
- [快速开始](#快速开始)
- [代码结构](#代码结构)
- [功能详解](#功能详解)
- [配置说明](#配置说明)
- [使用示例](#使用示例)
- [学习要点](#学习要点)

---

## ✨ 功能概览

这个插件实现了所有类型的插件接口，包括：

### 1️⃣ 基础功能（必须）
- ✅ 插件元数据定义
- ✅ 配置管理

### 2️⃣ 生命周期管理
- ✅ `OnLoad` - 插件加载
- ✅ `OnUnload` - 插件卸载
- ✅ `OnEnable` - 插件启用
- ✅ `OnDisable` - 插件禁用

### 3️⃣ 事件处理
- ✅ 文本消息处理
- ✅ 图片消息处理
- ✅ 会话过期事件

### 4️⃣ 命令系统
- ✅ `/example echo` - 回显命令
- ✅ `/example status` - 状态查询
- ✅ `/example config` - 配置修改
- ✅ `/example session` - 会话查看

### 5️⃣ 被调用能力
- ✅ `format` - 文本格式化
- ✅ `validate` - 长度验证
- ✅ `stats` - 统计信息

### 6️⃣ 高级功能
- ✅ 会话管理和劫持
- ✅ 线程安全的状态管理
- ✅ 统计和监控
- ✅ 配置热更新

---

## 🚀 快速开始

### 1. 构建插件

```bash
cd plugins/example
go mod tidy
go build -o example.exe
```

### 2. 配置插件

在 `plugins/config.toml` 中添加：

```toml
[example]
enable = true
mode = 'blacklist'
limits = []

[example.config]
echo_prefix = "[回显] "
reply_enabled = true
max_length = 500
```

### 3. 运行

将编译好的 `example.exe` 放在 `plugins/` 目录下，启动 Host 即可。

---

## 📁 代码结构

```go
plugins/example/
├── main.go           # 主文件（本文档讲解的核心）
├── go.mod           # Go 模块定义
└── README.md        # 本文档
```

### 主要组件

```go
// 1. 配置结构
type Config struct {
    EchoPrefix   string
    ReplyEnabled bool
    MaxLength    int
}

// 2. 插件主结构
type ExamplePlugin struct {
    plugin.ConfigAbility[Config]  // 配置能力
    message message.Ability        // 消息能力
    contact contact.Ability        // 联系人能力
    session plugin.SessionAbility  // 会话能力
    caller  plugin.CallerAbility   // 调用能力
    
    // 内部状态
    sessions map[string]*SessionData
    stats    Stats
}

// 3. 命令定义
type EchoCommand struct { ... }
type StatusCommand struct { ... }
type ConfigCommand struct { ... }
type SessionCommand struct { ... }
```

---

## 🔍 功能详解

### 1. 基础插件接口

#### GetMetadata() - 插件元数据

```go
func (p *ExamplePlugin) GetMetadata() *plugin.Metadata {
    return &plugin.Metadata{
        Name:        "example",      // 唯一标识
        Author:      "Golem Team",
        Version:     "2.0.0",
        Description: "完整示例插件",
        Priority:    0,              // 优先级
        Next:        false,          // 是否传递
        AlwaysRun:   false,          // 是否总运行
    }
}
```

**学习要点**：
- `Name` 是插件的唯一标识，对应配置文件中的段名
- `Priority` 控制执行顺序，数值越小越先执行
- `Next: true` 表示成功处理事件后，允许后续插件继续处理
- `AlwaysRun: true` 表示即使会话被劫持也继续响应事件

---

### 2. 生命周期管理

#### OnLoad() - 插件加载

```go
func (p *ExamplePlugin) OnLoad() error {
    // 初始化内部状态
    p.sessions = make(map[string]*SessionData)
    p.stats = Stats{}
    
    // 设置默认配置
    if p.Config.EchoPrefix == "" {
        p.Config.EchoPrefix = "[Echo] "
    }
    
    return nil
}
```

**使用场景**：
- 初始化数据结构
- 加载持久化数据
- 设置默认配置
- 建立外部连接

#### OnUnload() - 插件卸载

```go
func (p *ExamplePlugin) OnUnload() error {
    // 输出统计信息
    slog.Info("统计", "messages", p.stats.TotalMessages)
    
    // 清理资源
    p.sessions = nil
    
    return nil
}
```

**使用场景**：
- 保存状态
- 关闭连接
- 释放资源
- 输出统计

---

### 3. 事件处理

#### GetSubscriptions() - 事件订阅

```go
func (p *ExamplePlugin) GetSubscriptions() []string {
    return []string{
        message.TypeText.Topic(),    // 推荐：使用类型常量
        message.TypeImage.Topic(),   // 图片消息
        "session::expired",          // 或使用字符串
    }
}
```

**推荐做法**：使用 `message.TypeXxx.Topic()` 避免手写字符串出错。

**可用事件类型**：
- `message.TypeText.Topic()` - 文本消息
- `message.TypeImage.Topic()` - 图片消息
- `message::video` - 视频消息
- `message::voice` - 语音消息
- `message::file` - 文件消息
- `contact::add` - 好友添加
- `chatroom::join` - 入群事件
- `session::expired` - 会话过期

#### OnEvent() - 事件处理

```go
func (p *ExamplePlugin) OnEvent(e *plugin.Event) (bool, error) {
    // 返回值说明：
    // bool: true 表示已处理，false 表示未处理
    // error: 错误信息，nil 表示无错误
    
    switch e.Topic {
    case "message::text":
        return p.handleTextMessage(e)
    case "session::expired":
        return p.handleSessionExpired(e)
    }
    
    return false, nil
}
```

**处理文本消息示例**：

```go
func (p *ExamplePlugin) handleTextMessage(e *plugin.Event) (bool, error) {
    // 1. 解析消息
    msg := e.Payload.(*plugin.Event_Message)
    content := msg.Message.Content
    
    // 2. 劫持会话（10秒内此用户的消息只给这个插件）
    p.session.Hold(p, e.Sender, 10*time.Second)
    
    // 3. 发送回复
    _, err := p.message.Send(&message.Message{
        Type:     message.TypeText,  // 必须指定消息类型
        Receiver: &contact.Contact{Username: e.Sender},
        Content:  p.Config.EchoPrefix + content,
    })
    
    return true, err
}
```

**学习要点**：
- 使用 `e.Payload` 获取具体的事件数据
- 返回 `true` 表示已处理，其他插件将不再收到此事件
- 返回 `false` 表示未处理，事件继续传递
- `session.Hold()` 可以劫持会话，实现多轮对话

---

### 4. 命令系统

#### 定义命令

```go
type EchoCommand struct {
    // 命令元数据（必须是第一个字段，类型为空结构体）
    _ struct{} `cmd:"example echo" help:"回显文本" usage:"/example echo <text>" example:"/example echo hello"`
    
    // 命令参数
    Text   string `arg:"text" help:"回显内容" required:"true" variadic:"true"`
    Upper  bool   `flag:"upper" help:"转换为大写"`
    Prefix string `flag:"prefix" help:"自定义前缀"`
}
```

**标签说明**：
- `cmd`: 命令名称
- `help`: 命令说明
- `usage`: 使用方法
- `example`: 示例
- `arg`: 位置参数
- `flag`: 可选标志
- `required`: 是否必填
- `variadic`: 是否接受多个值

#### 注册和处理命令

```go
// 1. 实现接口
func (p *ExamplePlugin) GetCommands() []string {
    return plugin.CommandCommands()  // 自动返回已注册的命令
}

func (p *ExamplePlugin) OnCommand(cmd *plugin.Command) (string, error) {
    return plugin.DispatchCommand(cmd)  // 自动分发到处理函数
}

// 2. 定义处理函数
func (p *ExamplePlugin) handleEcho(cmd EchoCommand) (string, error) {
    text := cmd.Text
    if cmd.Upper {
        text = strings.ToUpper(text)
    }
    return cmd.Prefix + text, nil
}

// 3. 在 main 中注册
func main() {
    p := &ExamplePlugin{}
    plugin.RegisterCommand(p.handleEcho)
    plugin.Start(p)
}
```

**使用示例**：
```
/example echo hello                    # 基本用法
/example echo hello --upper            # 转大写
/example echo hello --prefix ">> "     # 自定义前缀
/example status                        # 查看状态
/example config prefix [新前缀]         # 修改配置
/example session                       # 查看会话
```

---

### 5. 被调用能力

允许其他插件调用本插件的功能：

```go
// 1. 声明提供的能力
func (p *ExamplePlugin) GetCapabilities() []string {
    return []string{
        "example:format",    // 使用前缀确保唯一性
        "example:validate",
        "example:stats",
    }
}

// 2. 处理调用
func (p *ExamplePlugin) OnCall(capability string, args map[string]string) (string, []byte, error) {
    switch capability {
    case "example:format":
        text := args["text"]
        result := p.Config.EchoPrefix + text
        return "text", []byte(result), nil
        
    case "example:validate":
        text := args["text"]
        valid := len(text) <= p.Config.MaxLength
        return "bool", []byte(fmt.Sprint(valid)), nil
        
    case "example:stats":
        stats := fmt.Sprintf(`{"messages":%d}`, p.stats.TotalMessages)
        return "json", []byte(stats), nil
    }
    
    return "", nil, fmt.Errorf("不支持的能力: %s", capability)
}
```

**其他插件调用方式**：

```go
// 在其他插件中
result, data, err := p.caller.CallPlugin("example", "example:format", map[string]string{
    "text": "hello",
})
// result: "text"
// data: []byte("[Echo] hello")
```

---

### 6. 高级功能

#### 会话管理

```go
type SessionData struct {
    LastMessage string
    Count       int
    StartTime   time.Time
}

// 更新会话
func (p *ExamplePlugin) updateSession(sender, content string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if session, ok := p.sessions[sender]; ok {
        session.LastMessage = content
        session.Count++
    } else {
        p.sessions[sender] = &SessionData{
            LastMessage: content,
            Count:       1,
            StartTime:   time.Now(),
        }
    }
}

// 劫持会话
p.session.Hold(p, e.Sender, 10*time.Second)
```

#### 线程安全

```go
// 使用读写锁保护共享状态
type ExamplePlugin struct {
    mu       sync.RWMutex
    sessions map[string]*SessionData
}

// 读操作
func (p *ExamplePlugin) getSession(sender string) *SessionData {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.sessions[sender]
}

// 写操作
func (p *ExamplePlugin) setSession(sender string, data *SessionData) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.sessions[sender] = data
}
```

#### 配置热更新

```go
// 修改配置
p.Config.EchoPrefix = "新前缀"

// 保存到文件
if err := p.SaveConfig(p); err != nil {
    slog.Error("保存失败", "err", err)
}
```

---

## ⚙️ 配置说明

```toml
[example]
enable = true           # 是否启用插件
mode = 'blacklist'     # 过滤模式：whitelist/blacklist/disabled
limits = []            # 黑白名单列表

[example.config]
echo_prefix = "[回显] "     # 回显消息的前缀
reply_enabled = true       # 是否自动回复
max_length = 500          # 消息最大长度
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `echo_prefix` | string | `[Echo] ` | 回显消息的前缀 |
| `reply_enabled` | bool | `true` | 是否启用自动回复 |
| `max_length` | int | `500` | 允许的最大消息长度 |

---

## 💡 使用示例

### 场景 1：自动回显

**配置**：
```toml
[example.config]
reply_enabled = true
echo_prefix = "🔊 "
```

**对话**：
```
用户: hello
插件: 🔊 hello

用户: how are you?
插件: 🔊 how are you?
```

### 场景 2：命令使用

```
用户: /example echo hello world --upper
插件: [ECHO] HELLO WORLD

用户: /example status
插件: 📊 Example 插件状态
      版本: 2.0.0
      ...

用户: /example config prefix ">> "
插件: ✅ 配置已更新: prefix = >> 

用户: /example session
插件: 🔄 活跃会话列表:
      👤 user123
        - 消息数: 5
        ...
```

### 场景 3：被其他插件调用

```go
// 在其他插件中
result, data, err := caller.CallPlugin("example", "format", map[string]string{
    "text": "test message",
})
// 返回: "[Echo] test message"

valid, data, err := caller.CallPlugin("example", "validate", map[string]string{
    "text": "short",
})
// 返回: "true"
```

---

## 📖 学习要点

### 1. 必须实现的接口

```go
// 唯一必须实现的方法
func (p *ExamplePlugin) GetMetadata() *plugin.Metadata
```

### 2. 按需实现的接口

```go
// 生命周期
func (p *ExamplePlugin) OnLoad() error
func (p *ExamplePlugin) OnUnload() error

// 事件处理
func (p *ExamplePlugin) GetSubscriptions() []string
func (p *ExamplePlugin) OnEvent(*plugin.Event) (bool, error)

// 命令处理
func (p *ExamplePlugin) GetCommands() []string
func (p *ExamplePlugin) OnCommand(*plugin.Command) (string, error)

// 被调用
func (p *ExamplePlugin) GetCapabilities() []string
func (p *ExamplePlugin) OnCall(string, map[string]string) (string, []byte, error)
```

### 3. 能力注入

```go
type ExamplePlugin struct {
    // 声明需要的能力，Host 会自动注入
    message message.Ability
    contact contact.Ability
    session plugin.SessionAbility
    caller  plugin.CallerAbility
}
```

### 4. 配置管理

```go
type ExamplePlugin struct {
    // 使用 ConfigAbility 获得配置功能
    plugin.ConfigAbility[Config]
}

// 访问配置
prefix := p.Config.EchoPrefix

// 修改并保存
p.Config.EchoPrefix = "new prefix"
p.SaveConfig(p)
```

### 5. 最佳实践

1. **日志记录**：使用 `slog` 记录关键事件
2. **错误处理**：不要吞掉错误，返回给 Host
3. **线程安全**：使用 `sync.RWMutex` 保护共享状态
4. **资源清理**：在 `OnUnload` 中释放资源
5. **配置验证**：在 `OnLoad` 中设置默认值

---

## 🎯 下一步

1. **阅读代码**：仔细阅读 `main.go`，理解每个部分的作用
2. **修改配置**：尝试修改配置文件，观察效果
3. **添加功能**：尝试添加新的命令或事件处理
4. **创建插件**：基于这个示例创建自己的插件

---

## 📚 相关文档

- [插件开发指南](../readme.md)
- [SDK 使用指南](../../SDK_USAGE.md)
- [项目 README](../../readme.md)

---

**这是最完整的插件示例，涵盖了所有功能！祝你开发愉快！🎉**
