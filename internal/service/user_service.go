package service

import (
	"chat-room/internal/dao/pool"   // 引入数据库连接池
	"chat-room/internal/model"      // 引入数据模型包
	"chat-room/pkg/common/request"  // 引入通用请求包
	"chat-room/pkg/common/response" // 引入通用响应包
	"chat-room/pkg/errors"          // 引入自定义错误处理包
	"chat-room/pkg/global/log"      // 引入全局日志记录器
	"time"                          // 引入时间包，用于处理时间相关操作

	"github.com/google/uuid" // 引入UUID库，用于生成唯一标识符
)

// userService 结构体实现了用户服务的相关逻辑
type userService struct {
}

// UserService 是全局的用户服务实例
var UserService = new(userService)

// Register 函数用于注册新用户
func (u *userService) Register(user *model.User) error {
	db := pool.GetDB() // 获取数据库连接实例
	var userCount int64
	db.Model(user).Where("username", user.Username).Count(&userCount) // 检查用户名是否已经存在
	if userCount > 0 {
		return errors.New("user already exists") // 如果用户名已存在，返回错误
	}
	user.Uuid = uuid.New().String() // 生成新用户的UUID
	user.CreateAt = time.Now()      // 设置用户创建时间
	user.DeleteAt = 0               // 初始化删除时间为0

	db.Create(&user) // 保存新用户信息到数据库
	return nil
}

// Login 函数用于用户登录验证
func (u *userService) Login(user *model.User) bool {
	pool.GetDB().AutoMigrate(&user) // 自动迁移用户表结构
	log.Logger.Debug("user", log.Any("user in service", user))
	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "username = ?", user.Username) // 根据用户名查询用户信息
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))

	user.Uuid = queryUser.Uuid // 将查询到的用户UUID赋值给传入的user对象

	return queryUser.Password == user.Password // 返回密码匹配结果
}

// ModifyUserInfo 函数用于修改用户信息
func (u *userService) ModifyUserInfo(user *model.User) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "username = ?", user.Username) // 根据用户名查询用户信息
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id { // 如果用户不存在，返回错误
		return errors.New("用户不存在")
	}
	// 更新用户信息
	queryUser.Nickname = user.Nickname
	queryUser.Email = user.Email
	queryUser.Password = user.Password

	db.Save(queryUser) // 保存更新后的用户信息
	return nil
}

// GetUserDetails 函数根据用户UUID获取用户详细信息
func (u *userService) GetUserDetails(uuid string) model.User {
	var queryUser *model.User
	db := pool.GetDB()
	// 查询并返回用户的基本信息
	db.Select("uuid", "username", "nickname", "avatar").First(&queryUser, "uuid = ?", uuid)
	return *queryUser
}

// GetUserOrGroupByName 函数通过名称查找群组或者用户
func (u *userService) GetUserOrGroupByName(name string) response.SearchResponse {
	var queryUser *model.User
	db := pool.GetDB()
	db.Select("uuid", "username", "nickname", "avatar").First(&queryUser, "username = ?", name) // 查询用户信息

	var queryGroup *model.Group
	db.Select("uuid", "name").First(&queryGroup, "name = ?", name) // 查询群组信息

	// 将查询结果封装到响应结构中返回
	search := response.SearchResponse{
		User:  *queryUser,
		Group: *queryGroup,
	}
	return search
}

// GetUserList 函数根据用户UUID获取用户的好友列表
func (u *userService) GetUserList(uuid string) []model.User {
	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "uuid = ?", uuid) // 根据UUID查询用户信息
	var nullId int32 = 0
	if nullId == queryUser.Id { // 如果用户不存在，返回nil
		return nil
	}

	var queryUsers []model.User
	// 查询用户的好友列表
	db.Raw("SELECT u.username, u.uuid, u.avatar FROM user_friends AS uf JOIN users AS u ON uf.friend_id = u.id WHERE uf.user_id = ?", queryUser.Id).Scan(&queryUsers)

	return queryUsers
}

// AddFriend 函数用于添加好友
func (u *userService) AddFriend(userFriendRequest *request.FriendRequest) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "uuid = ?", userFriendRequest.Uuid) // 根据UUID查询用户信息
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id { // 如果用户不存在，返回错误
		return errors.New("用户不存在")
	}

	var friend *model.User
	db.First(&friend, "username = ?", userFriendRequest.FriendUsername) // 根据用户名查询好友信息
	if nullId == friend.Id {                                            // 如果好友不存在，返回错误
		return errors.New("已添加该好友")
	}

	// 创建好友关系
	userFriend := model.UserFriend{
		UserId:   queryUser.Id,
		FriendId: friend.Id,
	}

	var userFriendQuery *model.UserFriend
	db.First(&userFriendQuery, "user_id = ? and friend_id = ?", queryUser.Id, friend.Id) // 检查好友关系是否已存在
	if userFriendQuery.ID != nullId {
		return errors.New("该用户已经是你好友")
	}

	db.AutoMigrate(&userFriend) // 自动迁移好友表结构
	db.Save(&userFriend)        // 保存好友关系
	log.Logger.Debug("userFriend", log.Any("userFriend", userFriend))

	return nil
}

// ModifyUserAvatar 函数用于修改用户头像
func (u *userService) ModifyUserAvatar(avatar string, userUuid string) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "uuid = ?", userUuid) // 根据用户UUID查询用户信息

	if NULL_ID == queryUser.Id { // 如果用户不存在，返回错误
		return errors.New("用户不存在")
	}

	db.Model(&queryUser).Update("avatar", avatar) // 更新用户的头像信息
	return nil
}
