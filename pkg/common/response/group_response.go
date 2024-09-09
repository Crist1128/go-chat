package response

import "time"

// GroupResponse 结构体用于封装群组信息的响应
type GroupResponse struct {
	Uuid      string    `json:"uuid"`     // 群组的UUID
	GroupId   int32     `json:"groupId"`  // 群组的ID
	CreatedAt time.Time `json:"createAt"` // 群组的创建时间
	Name      string    `json:"name"`     // 群组的名称
	Notice    string    `json:"notice"`   // 群组的公告
}
