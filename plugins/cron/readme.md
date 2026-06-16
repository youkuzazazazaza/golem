# Cron 插件

定时任务插件，支持根据 cron 表达式定时调用其他插件的能力，并将结果发送给指定的接收者。

## 功能特性

- **标准 cron 表达式**：支持标准的 5 位或 6 位 cron 表达式
- **能力调用**：可以调用任何插件提供的能力
- **多目标发送**：支持同时向多个接收者发送结果
- **参数传递**：支持向能力传递自定义参数
- **动态管理**：通过命令实时增删改查定时任务，无需重启
- **多格式支持**：支持文本、图片、语音、视频等多种消息类型

## 快速开始

### 示例 1：每天早上 9 点发送新闻

```bash
/cron add -c "0 9 * * *" -p news.today -t wxid_abc123
```

### 示例 2：每 5 分钟发送一次提醒

```bash
/cron add -c "*/5 * * * *" -p text.to.image -t wxid_abc123 -a "context=喝水啦 bg_color=#e3f2fd"
```

### 示例 3：每周一早上 8 点向多个群发送周报

```bash
/cron add -c "0 8 * * 1" -p report.weekly -t "123456@chatroom,789012@chatroom"
```

## 命令说明

### `/cron add`

创建定时任务。

**用法**：
```bash
/cron add -c <cron> -p <capability> -t <targets> [-a args]
```

**参数**：
- `-c, --cron`：cron 表达式（必需），包含空格时需要加引号
- `-p, --capability`：要调用的能力名称（必需）
- `-t, --targets`：接收者 username（必需），多个用英文逗号分隔
- `-a, --args`：调用参数（可选），支持 `key=value` 列表或 JSON 对象

**示例**：
```bash
# 基本用法
/cron add -c "0 9 * * *" -p news.today -t wxid_abc123

# 多个接收者
/cron add -c "0 9 * * *" -p news.today -t "wxid_abc123,wxid_def456"

# 带参数（key=value 格式）
/cron add -c "*/5 * * * *" -p text.to.image -t wxid_abc123 -a "context=hello bg_color=#fff"

# 带参数（JSON 格式）
/cron add -c "0 9 * * *" -p weather.forecast -t wxid_abc123 -a '{"city":"Beijing","days":"3"}'

# 包含空格的 cron 表达式需要加引号
/cron add -c "0 9 * * *" -p morning.greet -t wxid_abc123
```

### `/cron list`

列出所有定时任务。

**用法**：
```bash
/cron list
```

**示例输出**：
```
定时任务：
1. cron="0 9 * * *"
capability=news.today
targets=wxid_abc123
args=0

2. cron="*/5 * * * *"
capability=text.to.image
targets=wxid_abc123
args=2

3. cron="0 8 * * 1"
capability=report.weekly
targets=123456@chatroom,789012@chatroom
args=1
```

### `/cron update`

更新现有定时任务。

**用法**：
```bash
/cron update -i <id> [-c <cron>] [-p <capability>] [-t <targets>] [-a args]
```

**参数**：
- `-i, --id`：任务序号（必需），通过 `/cron list` 查看
- `-c, --cron`：新的 cron 表达式（可选）
- `-p, --capability`：新的能力名称（可选）
- `-t, --targets`：新的接收者列表（可选）
- `-a, --args`：新的参数（可选），传空字符串可清空参数

**示例**：
```bash
# 更新 cron 表达式
/cron update -i 1 -c "0 10 * * *"

# 更新能力
/cron update -i 1 -p news.morning

# 更新接收者
/cron update -i 1 -t "wxid_abc123,wxid_def456"

# 更新参数
/cron update -i 1 -a "city=Shanghai days=5"

# 清空参数
/cron update -i 1 -a ""

# 同时更新多个字段
/cron update -i 1 -c "0 10 * * *" -p news.morning -t wxid_abc123
```

### `/cron delete`

删除定时任务。

**用法**：
```bash
/cron delete -i <id>
```

**参数**：
- `-i, --id`：任务序号（必需），通过 `/cron list` 查看

**示例**：
```bash
/cron delete -i 1
```

## Cron 表达式说明

### 表达式格式

标准 5 位格式：
```
* * * * *
│ │ │ │ │
│ │ │ │ └── 星期几 (0-6，0 表示周日)
│ │ │ └──── 月份 (1-12)
│ │ └────── 日期 (1-31)
│ └──────── 小时 (0-23)
└────────── 分钟 (0-59)
```

可选的 6 位格式（包含秒）：
```
* * * * * *
│ │ │ │ │ │
│ │ │ │ │ └── 星期几 (0-6)
│ │ │ │ └──── 月份 (1-12)
│ │ │ └────── 日期 (1-31)
│ │ └──────── 小时 (0-23)
│ └────────── 分钟 (0-59)
└──────────── 秒 (0-59)
```

### 特殊字符

- `*`：匹配所有值
- `,`：列举多个值，如 `1,3,5`
- `-`：范围，如 `1-5`
- `/`：步长，如 `*/5` 表示每 5 个单位

### 常用表达式示例

```bash
# 每分钟
"* * * * *"

# 每小时
"0 * * * *"

# 每天早上 9 点
"0 9 * * *"

# 每天早上 9 点到下午 5 点的每个整点
"0 9-17 * * *"

# 每 5 分钟
"*/5 * * * *"

# 每 30 分钟
"*/30 * * * *"

# 每小时的第 15 和 45 分钟
"15,45 * * * *"

# 每周一早上 8 点
"0 8 * * 1"

# 每月 1 号早上 9 点
"0 9 1 * *"

# 每个工作日早上 9 点
"0 9 * * 1-5"

# 每周一和周五下午 6 点
"0 18 * * 1,5"

# 每天早上 8:30
"30 8 * * *"

# 每 2 小时
"0 */2 * * *"

# 每天 9 点、12 点、18 点
"0 9,12,18 * * *"
```

## 能力调用说明

### 能力名称格式

能力名称通常使用 `插件名.能力名` 格式，例如：
- `news.today`：新闻插件的今日新闻能力
- `weather.forecast`：天气插件的预报能力
- `text.to.image`：文字转图片插件的转换能力

具体可用的能力取决于已安装的插件。

### 参数格式

参数支持两种格式：

#### 1. key=value 格式

多个参数用空格分隔：
```bash
/cron add -c "0 9 * * *" -p weather.forecast -t wxid_abc123 -a "city=Beijing days=3"
```

#### 2. JSON 格式

使用 JSON 对象（需要用单引号包裹）：
```bash
/cron add -c "0 9 * * *" -p weather.forecast -t wxid_abc123 -a '{"city":"Beijing","days":"3"}'
```

### 特殊参数

- `receiver`：接收者 username，**不需要手动指定**，插件会自动添加

### 返回值类型

能力执行后会返回一个 MIME 类型和数据：

| MIME 类型 | 说明 | 消息类型 |
|----------|------|---------|
| `text`, `json` | 文本内容 | 文本消息 |
| `image` | 图片数据 | 图片消息 |
| `voice` | 语音数据 | 语音消息 |
| `video` | 视频数据 | 视频消息 |
| `none` | 无需发送消息 | 不发送 |

## 配置文件

配置文件位于 `~/.golem/plugins/cron.toml`，首次启动时自动生成。

**配置示例**：
```toml
[[jobs]]
cron = "0 9 * * *"
target = ["wxid_abc123"]
capability = "news.today"

[[jobs]]
cron = "*/5 * * * *"
target = ["wxid_abc123"]
capability = "text.to.image"

  [jobs.args]
  context = "hello"
  bg_color = "#fff"

[[jobs]]
cron = "0 8 * * 1"
target = ["123456@chatroom", "789012@chatroom"]
capability = "report.weekly"

  [jobs.args]
  format = "markdown"
```

## 实际应用场景

### 场景 1：每日新闻推送

每天早上 9 点自动发送今日新闻摘要：
```bash
/cron add -c "0 9 * * *" -p news.today -t wxid_abc123
```

### 场景 2：天气预报提醒

每天早上 7 点发送当天天气预报：
```bash
/cron add -c "0 7 * * *" -p weather.forecast -t wxid_abc123 -a "city=Beijing"
```

### 场景 3：定时喝水提醒

工作日每 2 小时提醒喝水：
```bash
/cron add -c "0 9-17/2 * * 1-5" -p text.to.image -t wxid_abc123 -a "context=该喝水了💧"
```

### 场景 4：周报生成

每周一早上 9 点生成并发送周报：
```bash
/cron add -c "0 9 * * 1" -p report.weekly -t "123456@chatroom,789012@chatroom"
```

### 场景 5：每月账单提醒

每月 1 号早上 9 点发送账单统计：
```bash
/cron add -c "0 9 1 * *" -p finance.bill -t wxid_abc123 -a "month=last"
```

### 场景 6：股票盘前提醒

工作日早上 9:15 发送股市开盘提醒：
```bash
/cron add -c "15 9 * * 1-5" -p stock.market -t wxid_abc123 -a "action=open"
```

### 场景 7：会议提醒

每周三下午 2:55 提醒周会：
```bash
/cron add -c "55 14 * * 3" -p reminder.meeting -t "123456@chatroom" -a "title=周例会 time=15:00"
```

### 场景 8：健康打卡

每天晚上 9 点提醒健康打卡：
```bash
/cron add -c "0 21 * * *" -p health.checkin -t wxid_abc123
```

## 工作流程

```
1. 到达预定时间
   ↓
2. Cron 插件被触发
   ↓
3. 调用指定的能力（传递参数）
   ↓
4. 能力返回 MIME 类型和数据
   ↓
5. 根据 MIME 类型构建消息
   ↓
6. 发送消息给接收者
```

## 注意事项

### 1. Cron 表达式

- **包含空格必须加引号**：`/cron add -c "0 9 * * *" ...`
- **时区**：使用服务器本地时区
- **最小间隔**：建议不低于 1 分钟，避免频繁调用
- **验证**：建议使用在线 cron 表达式工具验证

### 2. 能力调用

- **能力必须存在**：确保对应的插件已安装并启用
- **参数格式**：确保参数格式正确，否则能力可能无法正常工作
- **返回值**：确保能力返回正确的 MIME 类型和数据

### 3. 接收者

- **username 格式**：
  - 个人：`wxid_xxx` 或 `xxx@weixin`
  - 群聊：`xxxxx@chatroom`
- **多个接收者**：用英文逗号分隔，不要有空格（除非在引号内）
- **有效性**：确保接收者存在且可以接收消息

### 4. 性能考虑

- **并发任务**：多个任务可能同时执行，注意资源占用
- **大量接收者**：向大量接收者发送消息可能需要较长时间
- **失败处理**：某个接收者发送失败不会影响其他接收者

### 5. 配置管理

- **立即生效**：通过命令修改的任务立即生效
- **持久化**：任务配置自动保存到配置文件
- **重启恢复**：重启后自动加载配置文件中的任务

## 故障排查

### 1. 任务没有执行

**检查清单**：
- 确认 cron 表达式正确（使用在线工具验证）
- 确认当前时间是否匹配表达式
- 查看日志是否有错误信息
- 确认插件已启用

### 2. 能力调用失败

**常见原因**：
- 能力名称错误或不存在
- 依赖的插件未安装或未启用
- 参数格式错误或缺少必需参数
- 能力内部执行出错

**解决方法**：
- 使用 `/cron list` 检查配置
- 查看日志获取详细错误信息
- 确认能力所在插件的状态
- 验证参数格式和内容

### 3. 消息发送失败

**常见原因**：
- 接收者 username 错误或不存在
- 接收者无法接收消息（被拉黑、删除好友等）
- 返回的 MIME 类型不支持
- 数据为空或格式错误

**解决方法**：
- 确认接收者 username 正确
- 测试是否可以手动发送消息给该接收者
- 检查能力返回的数据格式
- 查看日志获取详细错误信息

### 4. 查看日志

日志会记录以下信息：
- 任务创建成功/失败
- 任务执行时间
- 能力调用成功/失败
- 消息发送成功/失败

日志示例：
```
[INFO] [cron] 定时任务创建成功 id=1 cron="0 9 * * *" capability=news.today
[INFO] [cron] 定时任务执行并发送结果成功 cron="0 9 * * *" capability=news.today target=wxid_abc123 mime=text
[ERROR] [cron] 定时任务调用能力失败 cron="0 9 * * *" capability=news.today target=wxid_abc123 err=能力不存在
```

## 最佳实践

### 1. 合理设置执行频率

- **高频任务**（分钟级）：确保能力执行快速，避免积压
- **低频任务**（小时/天级）：可以执行相对耗时的操作
- **避免尖峰**：多个任务不要集中在同一时刻

### 2. 错误处理

- **监控日志**：定期检查日志，及时发现问题
- **容错设计**：某个接收者失败不应影响其他接收者
- **重试机制**：能力内部应实现必要的重试逻辑

### 3. 参数设计

- **明确必需参数**：在能力文档中说明哪些参数是必需的
- **提供默认值**：对可选参数提供合理的默认值
- **参数验证**：能力应验证参数的有效性

### 4. 测试建议

- **先测试能力**：确保能力可以正常工作
- **使用短周期测试**：先用 `*/1 * * * *` 测试，确认无误后再改为实际周期
- **小范围测试**：先向自己发送，确认无误后再扩大接收者范围

## 开发信息

- **插件名称**：cron
- **版本**：0.0.1
- **作者**：golem
- **依赖**：robfig/cron v3.0.1
- **SDK 版本**：golem/sdk v0.1.1

## 扩展阅读

- [Cron 表达式在线生成器](https://crontab.guru/)
- [Go cron 库文档](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Golem 插件开发指南](../../readme.md)