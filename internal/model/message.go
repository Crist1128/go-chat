package model

import (
	"gorm.io/plugin/soft_delete" // 引入GORM的软删除插件，用于实现逻辑删除功能
	"time"                       // 引入时间包，用于处理时间相关操作
)

// Message 结构体表示消息的数据模型
type Message struct {
	ID          int32                 `json:"id" gorm:"primarykey"`                                                        // ID为主键，使用整型，自增
	CreatedAt   time.Time             `json:"createAt"`                                                                    // CreatedAt记录消息的创建时间
	UpdatedAt   time.Time             `json:"updatedAt"`                                                                   // UpdatedAt记录消息的最后更新时间
	DeletedAt   soft_delete.DeletedAt `json:"deletedAt"`                                                                   // DeletedAt用于软删除字段，标记记录是否被逻辑删除
	FromUserId  int32                 `json:"fromUserId" gorm:"index"`                                                     // FromUserId为发送消息的用户ID，数据库中为此字段创建索引
	ToUserId    int32                 `json:"toUserId" gorm:"index;comment:'发送给端的id，可为用户id或者群id'"`                         // ToUserId为接收消息的用户ID或群组ID，数据库中为此字段创建索引
	Content     string                `json:"content" gorm:"type:varchar(2500)"`                                           // Content存储消息内容，最大长度2500字符
	MessageType int16                 `json:"messageType" gorm:"comment:'消息类型：1单聊，2群聊'"`                                   // MessageType标识消息的类型，1表示单聊，2表示群聊
	ContentType int16                 `json:"contentType" gorm:"comment:'消息内容类型：1文字 2.普通文件 3.图片 4.音频 5.视频 6.语音聊天 7.视频聊天'"` // ContentType标识消息的内容类型，例如文字、文件、图片、音频等
	Pic         string                `json:"pic" gorm:"type:text;comment:'缩略图'"`                                          // Pic存储消息的缩略图地址，用于图片或视频的预览
	Url         string                `json:"url" gorm:"type:varchar(350);comment:'文件或者图片地址'"`                             // Url存储消息内容的URL，例如文件或图片的存储地址
}
