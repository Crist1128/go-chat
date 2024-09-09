package config

import (
	"fmt" // 引入格式化输入输出的标准库

	"github.com/spf13/viper" // 引入Viper库，用于配置管理
)

// TomlConfig 结构体表示项目的总体配置
type TomlConfig struct {
	AppName        string         // 应用程序名称
	MySQL          MySQLConfig    // MySQL数据库配置
	Log            LogConfig      // 日志配置
	StaticPath     PathConfig     // 静态文件路径配置
	MsgChannelType MsgChannelType // 消息队列类型及相关配置
}

// MySQLConfig 结构体表示MySQL相关配置
type MySQLConfig struct {
	Host        string // 数据库主机地址
	Name        string // 数据库名称
	Password    string // 数据库密码
	Port        int    // 数据库端口号
	TablePrefix string // 数据库表前缀
	User        string // 数据库用户名
}

// LogConfig 结构体表示日志配置
type LogConfig struct {
	Path  string // 日志文件保存路径
	Level string // 日志级别（如INFO、DEBUG、ERROR等）
}

// PathConfig 结构体表示静态文件路径的配置
type PathConfig struct {
	FilePath string // 静态文件路径
}

// MsgChannelType 结构体表示消息队列类型及其相关配置信息
// 如果使用Go的channel，则为单机使用；如果使用Kafka，则支持分布式扩展
type MsgChannelType struct {
	ChannelType string // 消息通道类型（Gochannel或Kafka）
	KafkaHosts  string // Kafka的主机地址列表
	KafkaTopic  string // Kafka的主题
}

// c 是一个TomlConfig类型的全局变量，用于存储读取到的配置信息
var c TomlConfig

// init 函数在包被引入时自动执行，用于初始化配置
func init() {
	// 设置配置文件的文件名
	viper.SetConfigName("config")
	// 设置配置文件的类型为toml
	viper.SetConfigType("toml")
	// 设置配置文件的查找路径，这里设置为当前目录
	viper.AddConfigPath(".")
	// 自动识别环境变量并覆盖配置
	viper.AutomaticEnv()
	// 读取配置文件内容
	err := viper.ReadInConfig()
	if err != nil {
		// 如果读取配置文件失败，程序终止并抛出错误
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 将读取到的配置文件内容解析到c变量中
	viper.Unmarshal(&c)
}

// GetConfig 函数用于获取全局的配置实例
func GetConfig() TomlConfig {
	return c
}
