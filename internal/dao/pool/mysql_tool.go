package pool

import (
	"chat-room/config"     // 引入配置包，用于读取数据库配置信息
	"fmt"                  // 引入fmt包，用于格式化字符串
	"gorm.io/driver/mysql" // 引入GORM的MySQL驱动
	"gorm.io/gorm"         // 引入GORM包，用于ORM操作
	"gorm.io/gorm/logger"  // 引入GORM的日志模块，用于记录数据库操作日志
)

var _db *gorm.DB // 定义一个全局变量，存储数据库连接实例

// init 函数在包被初始化时自动执行，负责数据库连接的初始化
func init() {
	// 从配置文件中读取数据库连接信息
	username := config.GetConfig().MySQL.User     // 数据库用户名
	password := config.GetConfig().MySQL.Password // 数据库密码
	host := config.GetConfig().MySQL.Host         // 数据库主机地址
	port := config.GetConfig().MySQL.Port         // 数据库端口号
	Dbname := config.GetConfig().MySQL.Name       // 数据库名称
	timeout := "10s"                              // 连接超时，设置为10秒

	// 拼接DSN（数据源名称），用于MySQL连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)

	var err error

	// 使用GORM库连接MySQL数据库，并设置日志级别为Info
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置GORM的日志模式为Info级别
	})
	if err != nil {
		// 如果连接数据库失败，终止程序并输出错误信息
		panic("连接数据库失败, error=" + err.Error())
	}

	// 获取底层的sql.DB实例，用于设置连接池参数
	sqlDB, _ := _db.DB()

	// 设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置数据库连接池最大连接数为100
	sqlDB.SetMaxIdleConns(20)  // 设置连接池最大空闲连接数为20
}

// GetDB 函数用于返回全局的数据库连接实例
func GetDB() *gorm.DB {
	return _db
}
