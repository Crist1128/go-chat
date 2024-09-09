package response

import "chat-room/internal/model"

// SearchResponse 结构体用于封装用户或群组的查询响应
type SearchResponse struct {
	User  model.User  `json:"user"`  // 查询到的用户信息
	Group model.Group `json:"group"` // 查询到的群组信息
}
