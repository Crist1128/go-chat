package response

// ResponseMsg 结构体用于封装通用的响应消息
type ResponseMsg struct {
	Code int         `json:"code"` // 响应码，0表示成功，-1表示失败
	Msg  string      `json:"msg"`  // 响应消息
	Data interface{} `json:"data"` // 响应数据，可以是任意类型
}

// SuccessMsg 函数返回一个表示成功的响应消息
func SuccessMsg(data interface{}) *ResponseMsg {
	msg := &ResponseMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: data,
	}
	return msg
}

// FailMsg 函数返回一个表示失败的响应消息
func FailMsg(msg string) *ResponseMsg {
	msgObj := &ResponseMsg{
		Code: -1,
		Msg:  msg,
	}
	return msgObj
}

// FailCodeMsg 函数返回一个自定义响应码的失败消息
func FailCodeMsg(code int, msg string) *ResponseMsg {
	msgObj := &ResponseMsg{
		Code: code,
		Msg:  msg,
	}
	return msgObj
}
