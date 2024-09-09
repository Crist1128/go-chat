package model

import (
	"gorm.io/plugin/soft_delete" // 引入GORM的软删除插件，用于实现逻辑删除功能
	"time"                       // 引入时间包，用于处理时间相关操作
)

// UserFriend 结构体表示用户好友关系的数据模型
type UserFriend struct {
	ID        int32                 `json:"id" gorm:"primarykey"`                 // ID为主键，使用整型，自增
	CreatedAt time.Time             `json:"createAt"`                             // CreatedAt记录用户添加好友的时间
	UpdatedAt time.Time             `json:"updatedAt"`                            // UpdatedAt记录用户好友关系的最后更新时间
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`                            // DeletedAt用于软删除字段，标记记录是否被逻辑删除
	UserId    int32                 `json:"userId" gorm:"index;comment:'用户ID'"`   // UserId为用户的ID，数据库中为此字段创建索引
	FriendId  int32                 `json:"friendId" gorm:"index;comment:'好友ID'"` // FriendId为好友的ID，数据库中为此字段创建索引
}
