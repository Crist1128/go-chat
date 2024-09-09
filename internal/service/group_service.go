package service

import (
	"chat-room/internal/dao/pool"   // 引入数据库连接池
	"chat-room/internal/model"      // 引入数据模型包
	"chat-room/pkg/common/response" // 引入通用响应包，用于定义响应结构
	"chat-room/pkg/errors"          // 引入自定义错误包
	"github.com/google/uuid"        // 引入UUID库，用于生成唯一标识符
)

// groupService 结构体定义了群组服务的实现
type groupService struct {
}

// GroupService 是全局的群组服务实例，用于外部调用
var GroupService = new(groupService)

// GetGroups 函数根据用户的UUID获取该用户所属的群组列表
func (g *groupService) GetGroups(uuid string) ([]response.GroupResponse, error) {
	db := pool.GetDB() // 获取数据库连接实例

	// 自动迁移数据库表（如果表不存在则创建）
	migrate := &model.Group{}
	pool.GetDB().AutoMigrate(&migrate)
	migrate2 := &model.GroupMember{}
	pool.GetDB().AutoMigrate(&migrate2)

	var queryUser *model.User
	db.First(&queryUser, "uuid = ?", uuid) // 根据UUID查询用户信息

	if queryUser.Id <= 0 { // 如果用户不存在
		return nil, errors.New("用户不存在") // 返回错误
	}

	var groups []response.GroupResponse

	// 使用原生SQL查询用户所属的群组信息
	db.Raw("SELECT g.id AS group_id, g.uuid, g.created_at, g.name, g.notice FROM group_members AS gm LEFT JOIN `groups` AS g ON gm.group_id = g.id WHERE gm.user_id = ?",
		queryUser.Id).Scan(&groups)

	return groups, nil // 返回群组列表
}

// SaveGroup 函数用于保存新的群组信息，并将创建者加入到群组中
func (g *groupService) SaveGroup(userUuid string, group model.Group) {
	db := pool.GetDB() // 获取数据库连接实例
	var fromUser model.User
	db.Find(&fromUser, "uuid = ?", userUuid) // 根据用户UUID查询用户信息
	if fromUser.Id <= 0 {
		return // 如果用户不存在，直接返回
	}

	// 设置群组的创建者ID和UUID
	group.UserId = fromUser.Id
	group.Uuid = uuid.New().String()
	db.Save(&group) // 保存群组信息

	// 创建群组成员记录，将创建者加入群组
	groupMember := model.GroupMember{
		UserId:   fromUser.Id,
		GroupId:  group.ID,
		Nickname: fromUser.Username,
		Mute:     0,
	}
	db.Save(&groupMember) // 保存群组成员信息
}

// GetUserIdByGroupUuid 函数根据群组的UUID获取群组内的用户列表
func (g *groupService) GetUserIdByGroupUuid(groupUuid string) []model.User {
	var group model.Group
	db := pool.GetDB()                      // 获取数据库连接实例
	db.First(&group, "uuid = ?", groupUuid) // 根据群组UUID查询群组信息
	if group.ID <= 0 {
		return nil // 如果群组不存在，返回空列表
	}

	var users []model.User
	// 使用原生SQL查询群组成员信息
	db.Raw("SELECT u.uuid, u.avatar, u.username FROM `groups` AS g JOIN group_members AS gm ON gm.group_id = g.id JOIN users AS u ON u.id = gm.user_id WHERE g.id = ?",
		group.ID).Scan(&users)
	return users // 返回用户列表
}

// JoinGroup 函数用于将用户加入到指定的群组
func (g *groupService) JoinGroup(groupUuid, userUuid string) error {
	var user model.User
	db := pool.GetDB()                    // 获取数据库连接实例
	db.First(&user, "uuid = ?", userUuid) // 根据用户UUID查询用户信息
	if user.Id <= 0 {
		return errors.New("用户不存在") // 如果用户不存在，返回错误
	}

	var group model.Group
	db.First(&group, "uuid = ?", groupUuid) // 根据群组UUID查询群组信息
	if group.ID <= 0 {
		return errors.New("群组不存在") // 如果群组不存在，返回错误
	}
	var groupMember model.GroupMember
	db.First(&groupMember, "user_id = ? and group_id = ?", user.Id, group.ID) // 检查用户是否已经是群组成员
	if groupMember.ID > 0 {
		return errors.New("已经加入该群组") // 如果已经加入群组，返回错误
	}
	nickname := user.Nickname
	if nickname == "" { // 如果用户没有设置昵称，则使用用户名作为昵称
		nickname = user.Username
	}
	// 创建新的群组成员记录
	groupMemberInsert := model.GroupMember{
		UserId:   user.Id,
		GroupId:  group.ID,
		Nickname: nickname,
		Mute:     0,
	}
	db.Save(&groupMemberInsert) // 保存群组成员信息

	return nil // 返回成功
}
