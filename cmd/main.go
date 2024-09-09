package main

import (
	"chat-room/config"              // 引入配置包，用于加载和访问配置信息
	"chat-room/internal/kafka"      // 引入Kafka包，用于处理Kafka消息队列
	"chat-room/internal/router"     // 引入路由包，用于定义HTTP请求路由
	"chat-room/internal/server"     // 引入服务器包，用于管理WebSocket服务器
	"chat-room/pkg/common/constant" // 引入常量包，定义了项目中使用的常量
	"chat-room/pkg/global/log"      // 引入日志包，用于日志记录
	"net/http"                      // 提供HTTP服务器的功能
	"time"                          // 提供时间相关的功能
)

func main() {
	// 初始化日志系统，配置日志路径和日志级别
	log.InitLogger(config.GetConfig().Log.Path, config.GetConfig().Log.Level)
	// 记录当前的配置信息
	log.Logger.Info("config", log.Any("config", config.GetConfig()))

	// 如果消息通道类型是Kafka，初始化Kafka生产者和消费者
	if config.GetConfig().MsgChannelType.ChannelType == constant.KAFKA {
		// 初始化Kafka生产者，指定Kafka主题和主机地址
		kafka.InitProducer(config.GetConfig().MsgChannelType.KafkaTopic, config.GetConfig().MsgChannelType.KafkaHosts)
		// 初始化Kafka消费者，指定主机地址
		kafka.InitConsumer(config.GetConfig().MsgChannelType.KafkaHosts)
		// 启动一个goroutine来异步消费Kafka消息，并处理消息
		go kafka.ConsumerMsg(server.ConsumerKafkaMsg)
	}

	// 记录服务器启动信息
	log.Logger.Info("start server", log.String("start", "start web server..."))

	// 初始化路由
	newRouter := router.NewRouter()

	// 启动WebSocket服务器
	go server.MyServer.Start()

	// 配置并启动HTTP服务器
	s := &http.Server{
		Addr:           "0.0.0.0:8888",   // 监听所有网络接口上的8888端口
		Handler:        newRouter,        // 使用初始化的路由器处理HTTP请求
		ReadTimeout:    10 * time.Second, // 设置读取超时时间为10秒
		WriteTimeout:   10 * time.Second, // 设置写入超时时间为10秒
		MaxHeaderBytes: 1 << 20,          // 设置请求头的最大字节数为1MB
	}
	// 启动HTTP服务器并监听端口
	err := s.ListenAndServe()
	if nil != err {
		// 如果服务器启动失败，记录错误信息
		log.Logger.Error("server error", log.Any("serverError", err))
	}
}
