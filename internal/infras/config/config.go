package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxLifetime     time.Duration `mapstructure:"max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	PingTimeout     time.Duration `mapstructure:"ping_timeout"`
	PrepareStmt     bool          `mapstructure:"prepare_stmt"`
	SkipDefaultTx   bool          `mapstructure:"skip_default_transaction"`
}

// DSN 生成 MySQL 连接字符串
func (c *MySQLConfig) DSN() string {
	return c.Username + ":" + c.Password +
		"@tcp(" + c.Host + ":" + itoa(c.Port) + ")/" + c.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local" +
		"&timeout=" + c.DialTimeout.String() +
		"&readTimeout=" + c.ReadTimeout.String() +
		"&writeTimeout=" + c.WriteTimeout.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// GetMySQLConfig 获取 MySQL 配置
func GetMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		Host:            viper.GetString("mysql.host"),
		Port:            viper.GetInt("mysql.port"),
		Username:        viper.GetString("mysql.username"),
		Password:        viper.GetString("mysql.password"),
		Database:        viper.GetString("mysql.database"),
		MaxIdleConns:    viper.GetInt("mysql.max_idle_conns"),
		MaxOpenConns:    viper.GetInt("mysql.max_open_conns"),
		MaxLifetime:     viper.GetDuration("mysql.max_lifetime") * time.Second,
		ConnMaxIdleTime: viper.GetDuration("mysql.conn_max_idle_time") * time.Second,
		DialTimeout:     viper.GetDuration("mysql.dial_timeout") * time.Second,
		ReadTimeout:     viper.GetDuration("mysql.read_timeout") * time.Second,
		WriteTimeout:    viper.GetDuration("mysql.write_timeout") * time.Second,
		PingTimeout:     viper.GetDuration("mysql.ping_timeout") * time.Second,
		PrepareStmt:     viper.GetBool("mysql.prepare_stmt"),
		SkipDefaultTx:   viper.GetBool("mysql.skip_default_transaction"),
	}
}

// QueueConfig 消息队列配置
type QueueConfig struct {
	Driver           string           `mapstructure:"driver"`
	ConsumerGroup    string           `mapstructure:"consumer_group"`
	TopicConcurrency map[string]int   `mapstructure:"topic_concurrency"`
	DefaultMaxLen    int64            `mapstructure:"default_max_len"`
	TopicMaxLen      map[string]int64 `mapstructure:"topic_max_len"`
	TrimInterval     int              `mapstructure:"trim_interval"` // 秒
}

// GetQueueConfig 获取消息队列配置
func GetQueueConfig() *QueueConfig {
	return &QueueConfig{
		Driver:           viper.GetString("queue.driver"),
		ConsumerGroup:    viper.GetString("queue.consumer_group"),
		TopicConcurrency: getTopicConcurrency(),
		DefaultMaxLen:    viper.GetInt64("queue.default_max_len"),
		TopicMaxLen:      getTopicMaxLen(),
		TrimInterval:     viper.GetInt("queue.trim_interval"),
	}
}

func getTopicConcurrency() map[string]int {
	raw := viper.GetStringMap("queue.topic_concurrency")
	result := make(map[string]int, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case int:
			result[k] = val
		case float64:
			result[k] = int(val)
		}
	}
	return result
}

func getTopicMaxLen() map[string]int64 {
	raw := viper.GetStringMap("queue.topic_max_len")
	result := make(map[string]int64, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case int:
			result[k] = int64(val)
		case float64:
			result[k] = int64(val)
		}
	}
	return result
}

// GetTopicMaxLen 获取指定 topic 的最大长度，未配置则使用全局默认值
func (c *QueueConfig) GetTopicMaxLen(topic string) int64 {
	if maxLen, ok := c.TopicMaxLen[topic]; ok && maxLen > 0 {
		return maxLen
	}
	return c.DefaultMaxLen
}

// GetTopicConcurrency 获取指定 topic 的并发数，未配置则使用默认值 1
func (c *QueueConfig) GetTopicConcurrency(topic string) int {
	if concurrency, ok := c.TopicConcurrency[topic]; ok && concurrency > 0 {
		return concurrency
	}
	return 1
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT 密钥
	Expire int    `mapstructure:"expire"` // 过期时间（秒）
}

// GetJWTConfig 获取 JWT 配置
func GetJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret: viper.GetString("jwt.secret"),
		Expire: viper.GetInt("jwt.expire"),
	}
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Mode       string `mapstructure:"mode"`        // dev / prod
	Level      string `mapstructure:"level"`       // debug / info / warn / error
	SQLLevel   string `mapstructure:"sql_level"`   // SQL 日志级别
	FilePath   string `mapstructure:"file_path"`   // 日志文件路径，空则只输出到控制台
	MaxSize    int    `mapstructure:"max_size"`    // 单文件最大 MB
	MaxBackups int    `mapstructure:"max_backups"` // 保留旧文件数
	MaxAge     int    `mapstructure:"max_age"`     // 保留天数
	Compress   bool   `mapstructure:"compress"`    // 压缩归档
}

// GetLoggerConfig 获取日志配置
func GetLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Mode:       viper.GetString("logger.mode"),
		Level:      viper.GetString("logger.level"),
		SQLLevel:   viper.GetString("logger.sql_level"),
		FilePath:   viper.GetString("logger.file_path"),
		MaxSize:    viper.GetInt("logger.max_size"),
		MaxBackups: viper.GetInt("logger.max_backups"),
		MaxAge:     viper.GetInt("logger.max_age"),
		Compress:   viper.GetBool("logger.compress"),
	}
}

func InitConfig() {
	viper.SetConfigName("config") // 使用 config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/infras/config") // 优先从 config 目录加载
	viper.AddConfigPath(".")                      // 也支持从根目录加载
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("mysql.host", "127.0.0.1")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.username", "root")
	viper.SetDefault("mysql.password", "123456")
	viper.SetDefault("mysql.database", "gin_admin")
	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.max_lifetime", 3600)
	viper.SetDefault("mysql.conn_max_idle_time", 600)
	viper.SetDefault("mysql.dial_timeout", 5)
	viper.SetDefault("mysql.read_timeout", 10)
	viper.SetDefault("mysql.write_timeout", 10)
	viper.SetDefault("mysql.ping_timeout", 3)
	viper.SetDefault("mysql.prepare_stmt", true)
	viper.SetDefault("mysql.skip_default_transaction", true)

	viper.SetDefault("server.port", 9501)

	viper.SetDefault("redis.addr", "127.0.0.1:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("queue.driver", "redis")
	viper.SetDefault("queue.consumer_group", "Gonio-group")
	viper.SetDefault("queue.default_max_len", 20)
	viper.SetDefault("queue.trim_interval", 3600)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found, using defaults and environment variables")
	}
}
