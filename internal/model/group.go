package model

import (
	"gorm.io/plugin/soft_delete" // 引入GORM的软删除插件，用于实现逻辑删除功能
	"time"                       // 引入时间包，用于时间处理
)

// Group 结构体表示群组的数据模型
type Group struct {
	ID        int32                 `json:"id" gorm:"primarykey"`                                                        // ID为主键，使用整型，自增
	Uuid      string                `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"` // Uuid是群组的唯一标识，使用字符串类型，长度150字符，不能为空，并且在数据库中唯一
	CreatedAt time.Time             `json:"createAt"`                                                                    // CreatedAt记录群组的创建时间
	UpdatedAt time.Time             `json:"updatedAt"`                                                                   // UpdatedAt记录群组的最后更新时间
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`                                                                   // DeletedAt用于软删除字段，实际不会从数据库中物理删除，而是标记为已删除
	UserId    int32                 `json:"userId" gorm:"index;comment:'群主ID'"`                                          // UserId是群主的ID，使用整型，并且在数据库中创建索引
	Name      string                `json:"name" gorm:"type:varchar(150);comment:'群名称'"`                                 // Name是群组的名称，使用字符串类型，长度150字符
	Notice    string                `json:"notice" gorm:"type:varchar(350);comment:'群公告'"`                               // Notice是群组的公告，使用字符串类型，长度350字符
}
