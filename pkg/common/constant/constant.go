package constant

// 定义一些常量，用于消息类型、内容类型和消息队列类型的管理

const (
	HEAT_BEAT = "heatbeat" // 心跳消息，用于维持长连接的存活状态
	PONG      = "pong"     // Pong消息，通常用于回应心跳包

	// 消息类型常量，用于区分消息是单聊还是群聊
	MESSAGE_TYPE_USER  = 1 // 单聊消息
	MESSAGE_TYPE_GROUP = 2 // 群聊消息

	// 消息内容类型常量，用于区分消息的内容类型
	TEXT         = 1 // 文字消息
	FILE         = 2 // 文件消息
	IMAGE        = 3 // 图片消息
	AUDIO        = 4 // 音频消息
	VIDEO        = 5 // 视频消息
	AUDIO_ONLINE = 6 // 在线语音聊天
	VIDEO_ONLINE = 7 // 在线视频聊天

	// 消息队列类型常量，用于区分使用的消息队列
	GO_CHANNEL = "gochannel" // 使用Go内置的channel作为消息队列
	KAFKA      = "kafka"     // 使用Kafka作为消息队列
)
