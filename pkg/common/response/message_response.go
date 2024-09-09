package response

import "time"

// MessageResponse 结构体用于封装消息信息的响应
type MessageResponse struct {
	ID           int32     `json:"id" gorm:"primarykey"`                            // 消息的ID
	FromUserId   int32     `json:"fromUserId" gorm:"index"`                         // 发送消息的用户ID
	ToUserId     int32     `json:"toUserId" gorm:"index"`                           // 接收消息的用户或群组ID
	Content      string    `json:"content" gorm:"type:varchar(2500)"`               // 消息内容
	ContentType  int16     `json:"contentType" gorm:"comment:'消息内容类型：1文字，2语音，3视频'"` // 消息内容类型
	CreatedAt    time.Time `json:"createAt"`                                        // 消息的创建时间
	FromUsername string    `json:"fromUsername"`                                    // 发送消息的用户名
	ToUsername   string    `json:"toUsername"`                                      // 接收消息的用户名（用于单聊）
	Avatar       string    `json:"avatar"`                                          // 发送消息用户的头像
	Url          string    `json:"url"`                                             // 消息中包含的URL（用于文件或多媒体消息）
}
