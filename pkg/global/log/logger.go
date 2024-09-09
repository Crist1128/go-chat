package log

import (
	"os" // 引入os包，用于文件操作

	"github.com/natefinch/lumberjack" // 引入lumberjack包，用于日志文件切割
	"go.uber.org/zap"                 // 引入zap包，用于高性能日志记录
	"go.uber.org/zap/zapcore"         // 引入zapcore包，用于配置日志核心组件
)

// 定义常用的zap.Field类型别名，方便在其他地方使用
type Field = zap.Field

var (
	Logger  *zap.Logger   // 定义全局日志实例
	String  = zap.String  // 定义zap.String的快捷方式
	Any     = zap.Any     // 定义zap.Any的快捷方式
	Int     = zap.Int     // 定义zap.Int的快捷方式
	Float32 = zap.Float32 // 定义zap.Float32的快捷方式
)

// InitLogger 初始化日志记录器
// logpath 日志文件路径
// loglevel 日志级别
func InitLogger(logpath string, loglevel string) {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    100,     // 每个日志文件最大100MB，默认 100M
		MaxBackups: 30,      // 保留30个备份，默认不限
		MaxAge:     7,       // 保留7天，默认不限
		Compress:   true,    // 启用日志压缩
	}
	write := zapcore.AddSync(&hook) // 将日志分割器与zap日志系统关联

	// 设置日志级别
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel // Debug级别，输出所有级别日志
	case "info":
		level = zap.InfoLevel // Info级别，输出Info以上级别日志
	case "error":
		level = zap.ErrorLevel // Error级别，只输出Error和以上级别日志
	case "warn":
		level = zap.WarnLevel // Warn级别，只输出Warn和Error级别日志
	default:
		level = zap.InfoLevel // 默认Info级别
	}

	// 配置日志的编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间字段的键名
		LevelKey:       "level",                        // 日志级别字段的键名
		NameKey:        "logger",                       // 日志记录器名称字段的键名
		CallerKey:      "linenum",                      // 调用者（行号）字段的键名
		MessageKey:     "msg",                          // 日志消息字段的键名
		StacktraceKey:  "stacktrace",                   // 堆栈跟踪字段的键名
		LineEnding:     zapcore.DefaultLineEnding,      // 行结尾的格式
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 日志级别小写编码
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间格式为ISO8601
		EncodeDuration: zapcore.SecondsDurationEncoder, // 持续时间以秒为单位
		EncodeCaller:   zapcore.FullCallerEncoder,      // 显示完整的调用路径
		EncodeName:     zapcore.FullNameEncoder,        // 显示完整的记录器名称
	}

	// 设置日志级别的原子操作
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	var writes = []zapcore.WriteSyncer{write}
	// 如果是开发环境，同时在控制台输出
	if level == zap.DebugLevel {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	// 创建日志核心组件
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 使用控制台编码器输出
		zapcore.NewMultiWriteSyncer(writes...),   // 同时写入文件和控制台
		level,                                    // 日志级别
	)

	// 设置日志记录器的附加选项
	caller := zap.AddCaller()                                   // 启用行号显示
	development := zap.Development()                            // 启用开发模式
	filed := zap.Fields(zap.String("application", "chat-room")) // 添加全局字段（如应用名称）

	// 构造日志记录器
	Logger = zap.New(core, caller, development, filed)
	Logger.Info("Logger init success") // 初始化成功日志
}
