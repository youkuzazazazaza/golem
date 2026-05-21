//go:build lib

package config

// Config SDK配置
type Config struct {
	Core    CoreConfig    `toml:"core" comment:"核心配置"`
	Server  ServerConfig  `toml:"server" comment:"Web服务配置"`
	Storage StorageConfig `toml:"storage" comment:"存储配置（路径固定为./data，不可配置）"`
	Device  DeviceConfig  `toml:"device" comment:"设备配置"`
	Log     LogConfig     `toml:"log" comment:"日志配置"`
}

// CoreConfig 核心配置
type CoreConfig struct {
	License             string `toml:"license" comment:"许可证"`
	PrintBanner         bool   `toml:"print_banner" comment:"是否打印banner"`
	IgnoreHistory       bool   `toml:"ignore_history" comment:"是否忽略历史消息，默认true"`
	IgnoreTimeout       int64  `toml:"ignore_timeout" comment:"历史消息超时时间（秒），默认10，即忽略10秒之前的消息"`
	RandomDelay         bool   `toml:"random_delay" comment:"是否启用随机延时，默认为true"`
	RandomDelayMin      int    `toml:"random_delay_min" comment:"随机延时最小值（毫秒）"`
	RandomDelayMax      int    `toml:"random_delay_max" comment:"随机延时最大值（毫秒）"`
	RandomDelayStrategy string `toml:"random_delay_strategy" comment:"随机延时策略，可选值：every（每个消息随机延时），pending（只有消息多排队时才加延时）"`
	QrcodeApi           string `toml:"qrcode_api,omitempty" comment:"二维码服务qpi"`
	RSAVersion          byte   `toml:"rsa_version,omitempty" comment:"rsa密钥版本"`
}

// ServerConfig Web服务配置
type ServerConfig struct {
	Enable bool   `toml:"enable" comment:"启用web服务"`
	Port   int    `toml:"port" comment:"服务端口"`
	Host   string `toml:"host" comment:"监听地址"`
	Token  string `toml:"token" comment:"认证token，只有携带token的请求才处理"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Encrypt    bool   `toml:"encrypt" comment:"是否加密存储"`
	EncryptKey string `toml:"encrypt_key" comment:"加密密钥（encrypt=true时必填）"`
}

// DeviceConfig 设备配置
type DeviceConfig struct {
	ID   string `toml:"id" comment:"设备ID（为空时自动生成）"`
	Name string `toml:"name" comment:"设备名称（为空时自动生成）"`
	Type string `toml:"type" comment:"设备类型：iphone, ipad, mac, android, android_pad, windows, car"`
	// 高级配置（可选，不输出到默认配置）
	OsVersion     string `toml:"os_version,omitempty" comment:"系统版本（高级配置）"`
	ClientVersion uint32 `toml:"client_version,omitempty" comment:"客户端版本号（高级配置）"`
	InitKey       string `toml:"init_key,omitempty" comment:"初始公钥hex（高级配置）"`

	KeyVersion byte `toml:"-" comment:"密钥版本：144, 145, 146, 147"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `toml:"level" comment:"日志级别：debug, info, warn, error"`
	Output     string `toml:"output" comment:"输出方式：console, file, both"`
	AddSource  bool   `toml:"add_source,omitempty" comment:"添加源地址"`
	Ansi       bool   `toml:"ansi" comment:"是否启用ansi颜色输出"`
	MaxSize    int    `toml:"max_size" comment:"单文件最大大小(MB)"`
	MaxAge     int    `toml:"max_age" comment:"保留天数"`
	MaxBackups int    `toml:"max_backups" comment:"保留文件数"`
	Compress   bool   `toml:"compress" comment:"是否压缩旧文件"`
}
