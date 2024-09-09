package service

import (
	"chat-room/internal/dao/pool"   // 引入数据库连接池
	"chat-room/internal/model"      // 引入数据模型包
	"chat-room/pkg/common/constant" // 引入全局常量
	"chat-room/pkg/common/request"  // 引入通用请求包
	"chat-room/pkg/common/response" // 引入通用响应包
	"chat-room/pkg/errors"          // 引入自定义错误处理包
	"chat-room/pkg/global/log"      // 引入全局日志记录器
	"chat-room/pkg/protocol"        // 引入消息协议包
	"gorm.io/gorm"                  // 引入GORM ORM库
)

// NULL_ID 定义了一个常量表示无效的ID
const NULL_ID int32 = 0

// messageService 结构体实现消息服务的相关逻辑
type messageService struct {
}

// MessageService 是全局的消息服务实例
var MessageService = new(messageService)

// GetMessages 函数根据请求参数获取消息列表
func (m *messageService) GetMessages(message request.MessageRequest) ([]response.MessageResponse, error) {
	db := pool.GetDB() // 获取数据库连接实例

	// 自动迁移消息表结构，确保表存在
	migrate := &model.Message{}
	pool.GetDB().AutoMigrate(&migrate)

	// 处理用户消息查询
	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var queryUser *model.User
		db.First(&queryUser, "uuid = ?", message.Uuid) // 根据UUID查询用户

		if NULL_ID == queryUser.Id { // 如果用户不存在
			return nil, errors.New("用户不存在")
		}

		var friend *model.User
		db.First(&friend, "username = ?", message.FriendUsername) // 根据用户名查询好友信息
		if NULL_ID == friend.Id {
			return nil, errors.New("用户不存在")
		}

		var messages []response.MessageResponse

		// 查询两个用户之间的消息
		db.Raw("SELECT m.id, m.from_user_id, m.to_user_id, m.content, m.content_type, m.url, m.created_at, u.username AS from_username, u.avatar, to_user.username AS to_username  FROM messages AS m LEFT JOIN users AS u ON m.from_user_id = u.id LEFT JOIN users AS to_user ON m.to_user_id = to_user.id WHERE from_user_id IN (?, ?) AND to_user_id IN (?, ?)",
			queryUser.Id, friend.Id, queryUser.Id, friend.Id).Scan(&messages)

		return messages, nil
	}

	// 处理群组消息查询
	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		messages, err := fetchGroupMessage(db, message.Uuid) // 调用辅助函数获取群组消息
		if err != nil {
			return nil, err
		}

		return messages, nil
	}

	return nil, errors.New("不支持查询类型") // 返回不支持的查询类型错误
}

// fetchGroupMessage 函数根据群组UUID获取群组消息
func fetchGroupMessage(db *gorm.DB, toUuid string) ([]response.MessageResponse, error) {
	var group model.Group
	db.First(&group, "uuid = ?", toUuid) // 根据UUID查询群组
	if group.ID <= 0 {
		return nil, errors.New("群组不存在")
	}

	var messages []response.MessageResponse

	// 查询群组内的消息
	db.Raw("SELECT m.id, m.from_user_id, m.to_user_id, m.content, m.content_type, m.url, m.created_at, u.username AS from_username, u.avatar FROM messages AS m LEFT JOIN users AS u ON m.from_user_id = u.id WHERE m.message_type = 2 AND m.to_user_id = ?",
		group.ID).Scan(&messages)

	return messages, nil
}

// SaveMessage 函数保存消息记录到数据库
func (m *messageService) SaveMessage(message protocol.Message) {
	db := pool.GetDB() // 获取数据库连接实例
	var fromUser model.User
	db.Find(&fromUser, "uuid = ?", message.From) // 根据消息发送者的UUID查询用户信息
	if NULL_ID == fromUser.Id {
		log.Logger.Error("SaveMessage not find from user", log.Any("SaveMessage not find from user", fromUser.Id))
		return
	}

	var toUserId int32 = 0

	// 处理单聊消息的保存
	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var toUser model.User
		db.Find(&toUser, "uuid = ?", message.To) // 根据消息接收者的UUID查询用户信息
		if NULL_ID == toUser.Id {
			return
		}
		toUserId = toUser.Id
	}

	// 处理群组消息的保存
	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var group model.Group
		db.Find(&group, "uuid = ?", message.To) // 根据消息接收者的UUID查询群组信息
		if NULL_ID == group.ID {
			return
		}
		toUserId = group.ID
	}

	// 创建并保存消息记录
	saveMessage := model.Message{
		FromUserId:  fromUser.Id,
		ToUserId:    toUserId,
		Content:     message.Content,
		ContentType: int16(message.ContentType),
		MessageType: int16(message.MessageType),
		Url:         message.Url,
	}
	db.Save(&saveMessage) // 保存消息到数据库
}
