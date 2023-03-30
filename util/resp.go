package util

import (
	"encoding/json"
	"fmt"
)

// ResMsg;http响应通用结构
type RespMsg struct {
	Code int
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 生成ResMsg对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// 将ResMsg对象转json的byte[]
func (resp *RespMsg) JSONBytes() []byte {
	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return b
}

// 将ResMsg对象转json格式的string
func (resp *RespMsg) JSONSring() string {
	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return string(b)
}
