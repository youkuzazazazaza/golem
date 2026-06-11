# Universal 插件说明

Universal 是一个规则驱动的通用 API 请求与消息发送插件。开发者可以通过命令或 TOML 配置添加规则，让插件在收到指定关键词消息后请求外部接口，并把接口结果以文本、图片、视频或表情发送回当前会话。

## 工作流程

1. 插件订阅 `message::text` 和 `message::app::quote` 事件。
2. 收到消息后，`parseIncomingText` 按第一个空格拆分文本：
   - 第一段是 `keyword`，用于匹配规则。
   - 剩余内容是 `param`，用于生成模板变量和 at 目标。
3. `ruleForKeyword` 使用 `keyword` 查找启用中的规则。
4. `extractQuote` 尝试解析引用消息，提取被引用者和引用文本。
5. `buildTemplateVars` 生成模板变量。
6. `executeRule` 渲染 URL、Headers、Body，发起 HTTP 请求，并按 `result_path` 提取结果。
7. 如果规则启用了 `continue_request`，会把第一次结果当作 URL 再请求一次。
8. `sendResult` 根据 `send_type` 构造消息并发送。
9. 如果规则 `at=true`，文本回复会在开头拼接 at 前缀，并在 `TextData.Reminds` 中传递对应 username。

## 规则配置

配置保存在插件配置 TOML 中，规则结构如下：

```toml
http_timeout_seconds = 15

[[rules]]
id = "cp"
keywords = ["cp"]
url = "https://api.example.com/cp?a={{arg_1}}&b={{arg_2}}"
method = "GET"
headers = ""
body = ""
send_type = "text"
result_path = "data.text"
at = false
continue_request = false
continue_method = "GET"
continue_result_path = ""
enabled = true
```

字段说明：

- `id`：规则唯一 ID。
- `keywords`：关键词列表。收到消息后只使用第一段 `keyword` 做全文匹配。
- `url`：请求地址模板，支持模板变量。
- `method`：HTTP 方法，默认 `GET`。
- `headers`：请求头，格式 `A=B;C=D`，支持模板变量。
- `body`：请求体模板，常用于 `POST` 或 `PUT`。
- `send_type`：发送类型，支持 `text`、`image`、`video`、`emoji`，默认 `text`。
- `result_path`：gjson 路径；为空时直接使用响应 body 原文。
- `at`：是否在文本回复中 at 参数或引用对应用户。
- `continue_request`：是否继续请求第一次结果中的 URL。
- `continue_method`：继续请求的 HTTP 方法，默认 `GET`。
- `continue_result_path`：继续请求响应的 gjson 路径；为空时直接使用响应 body 原文。
- `enabled`：是否启用；缺省时按启用处理。

## 命令用法

查看规则：

```text
/universal list
/universal get <id>
```

新增规则：

```text
/universal add <id> -k <keywords> -u <url> [options]
```

常用选项：

```text
-m, --method <method>
-H, --headers <headers>
-b, --body <body>
-t, --send-type <text|image|video|emoji>
-p, --result-path <gjson path>
--at <true|false>
-f, --continue
-M, --continue-method <method>
-P, --continue-result-path <gjson path>
```

更新规则：

```text
/universal update <id> [options]
```

清空字段：

```text
--clear-headers
--clear-body
--clear-result-path
--clear-continue-result-path
```

启用、禁用、删除：

```text
/universal enable <id>
/universal disable <id>
/universal delete <id>
```

## 模板变量

基础变量：

- `{{text}}`：完整消息文本。
- `{{keyword}}`：第一段关键词。
- `{{param}}`：关键词后的完整参数文本。
- `{{quote}}`：被引用文本。
- `{{quoter}}`：被引用者 display name。

位置参数：

- 普通参数会生成 `{{arg_1}}`、`{{arg_2}}`。
- 参数前缀 `@` 会在模板变量中移除。
- 如果存在引用消息，`{{quoter}}` 会自动成为 `{{arg_1}}`，文本参数整体后移。

命名参数：

```text
来张图 width=1920 height=1080
```

会生成：

```text
{{width}} = 1920
{{height}} = 1080
```

命名参数不会生成位置参数。

## At 回复规则

`at=true` 只对 `send_type=text` 有实际效果，因为当前消息协议只有 `TextData` 携带 `reminds` 字段。

At 目标解析顺序：

1. 引用消息目标优先。
2. 普通参数中的 `@用户` 目标随后追加。
3. 按 username 去重。

引用消息：

- username 使用 `refermsg.chatusr`。
- display name 使用 `refermsg.displayname`，为空时回退到 username。
- `refermsg.fromusr` 在群聊中是群 username，不能用于 at。

普通 @ 参数：

- 支持多个 `@用户`。
- nickname 从参数文本中的 `@用户` 提取。
- username 从原消息 `TextData.Reminds` 按顺序匹配。
- nickname 和 username 必须同时存在，否则跳过该目标。

示例：

```text
cp @张三 @李四
```

如果规则 `at=true`，回复文本会变成：

```text
@张三 @李四 接口返回结果
```

同时 `TextData.Reminds` 会传入张三和李四对应的 username。

## 示例

文本规则：

```text
/universal add cp -k cp -u "https://api.example.com/cp?a={{arg_1}}&b={{arg_2}}" -t text -p data.text --at true
```

触发：

```text
cp @张三 @李四
```

图片规则：

```text
/universal add pixiv -k 来张图 -u "https://api.example.com/image?width={{width}}&height={{height}}" -t image -p data.url
```

触发：

```text
来张图 width=1920 height=1080
```

关闭 at：

```text
/universal update cp --at false
```

继续请求：

```text
/universal add image -k 随机图 -u "https://api.example.com/random" -p data.url -t image -f -M GET
```

## 开发说明

主要文件：

- `main.go`：插件入口、事件处理和规则执行流程。
- `commands.go`：命令结构和命令处理。
- `rules.go`：规则索引、校验、默认值和展示。
- `incoming.go`：消息文本和引用消息解析。
- `template.go`：模板变量生成和模板渲染。
- `executor.go`：HTTP 请求、结果提取和继续请求。
- `send.go`：消息构造、媒体下载和发送。
- `mention.go`：at 目标收集和去重。
- `types.go`：规则、消息上下文和内部结构定义。

修改代码后建议执行：

```text
gofmt -w *.go
go test ./...
go build ./...
```

