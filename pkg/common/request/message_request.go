package request

// MessageRequest 结构体用于封装获取消息的请求参数
type MessageRequest struct {
	MessageType    int32  `json:"messageType"`    // 消息类型（单聊或群聊）
	Uuid           string `json:"uuid"`           // 当前用户的UUID
	FriendUsername string `json:"friendUsername"` // 好友的用户名（用于单聊时指定好友）
}
