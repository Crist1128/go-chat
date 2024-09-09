package model

import (
	"gorm.io/gorm" // 引入GORM包，用于ORM操作
	"time"         // 引入时间包，用于处理时间相关操作
)

// User 结构体表示用户的数据模型
type User struct {
	Id       int32      `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:'id'"`                                             // ID为主键，自增
	Uuid     string     `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`                   // UUID为用户的唯一标识符，长度150字符，不能为空，数据库中唯一索引
	Username string     `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'用户名'"`          // 用户名，必须字段，数据库中唯一且不能为空
	Password string     `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null; comment:'密码'"` // 密码，必须字段，长度150字符，不能为空
	Nickname string     `json:"nickname" gorm:"comment:'昵称'"`                                                                // 昵称，可选字段
	Avatar   string     `json:"avatar" gorm:"type:varchar(150);comment:'头像'"`                                                // 头像URL，长度150字符
	Email    string     `json:"email" gorm:"type:varchar(80);column:email;comment:'邮箱'"`                                     // 邮箱地址，长度80字符
	CreateAt time.Time  `json:"createAt"`                                                                                      // 创建时间，GORM 会自动填充该字段
	UpdateAt *time.Time `json:"updateAt"`                                                                                      // 更新时间，使用指针类型，允许在不更新时为nil
	DeleteAt int64      `json:"deleteAt"`                                                                                      // 逻辑删除时间，使用Unix时间戳格式存储
}

// BeforeUpdate 是 GORM 的一个钩子方法，会在更新操作之前执行
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdateAt", time.Now()) // 在更新操作之前自动更新 UpdateAt 字段为当前时间
	return nil
}
