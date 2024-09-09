package errors

// 定义一个自定义的错误类型，包含一个消息字段
type error struct {
	msg string // 错误信息
}

// Error 方法实现了 error 接口中的 Error() 方法，返回错误信息
func (e error) Error() string {
	return e.msg
}

// New 函数用于创建一个新的 error 实例，接受一个字符串作为错误信息
func New(msg string) error {
	return error{
		msg: msg, // 将传入的消息赋值给自定义错误的 msg 字段
	}
}
