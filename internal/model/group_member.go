package model

import (
	"gorm.io/plugin/soft_delete" // 引入GORM的软删除插件，用于实现逻辑删除功能
	"time"                       // 引入时间包，用于处理时间相关操作
)

// GroupMember 结构体表示群组成员的数据模型
type GroupMember struct {
	ID        int32                 `json:"id" gorm:"primarykey"`                           // ID为主键，使用整型，自增
	CreatedAt time.Time             `json:"createAt"`                                       // CreatedAt记录成员加入群组的时间
	UpdatedAt time.Time             `json:"updatedAt"`                                      // UpdatedAt记录成员最后一次信息更新的时间
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`                                      // DeletedAt用于软删除字段，标记记录是否被逻辑删除
	UserId    int32                 `json:"userId" gorm:"index;comment:'用户ID'"`             // UserId为用户ID，用于标识成员，数据库中创建索引
	GroupId   int32                 `json:"groupId" gorm:"index;comment:'群组ID'"`            // GroupId为群组ID，标识成员所属群组，数据库中创建索引
	Nickname  string                `json:"nickname" gorm:"type:varchar(350);comment:'昵称'"` // Nickname为成员在群中的昵称，使用字符串类型，最大长度350字符
	Mute      int16                 `json:"mute" gorm:"comment:'是否禁言'"`                     // Mute标识成员是否被禁言，0表示未禁言，1表示已禁言
}
