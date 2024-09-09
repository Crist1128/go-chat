package request

// FriendRequest 结构体用于封装添加好友的请求参数
type FriendRequest struct {
	Uuid           string // 发送请求的用户UUID
	FriendUsername string // 要添加的好友用户名
}
