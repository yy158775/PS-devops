package model

//消息传输的格式
type MessageResp struct {
	UserName string `json:"user_name"`
	Data     string `json:"data"`
}
