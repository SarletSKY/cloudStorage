package util

import (
	"encoding/json"
	"fmt"
	"log"
)

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 返回json对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// 返回[]byte
func (resp *RespMsg) JSONBytes() []byte {
	rBytes, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	return rBytes
}

// 返回string
func (resp *RespMsg) JSONString() string {
	rBytes, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	return string(rBytes)
}

//只返回code与msg
func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}

// 只反返回code与msg
func GetSimpleRespStreamString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}
