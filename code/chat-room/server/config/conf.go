package config

import (
	"encoding/json"
	"grpc-demo/server/logrus"
	"os"
)

type configuration struct {
	ServerInfo serverInfo
	RedisInfo  redisInfo
}

type redisInfo struct {
	RedisServerHost string
}

type serverInfo struct {
	RpcHost  string
	RespHost string
}

//这种大小写交错的模式，就是给你一个实体变量让你输出，其他都不让你输出
var Congfiguration configuration

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error": err,
		}).Fatal("Open config.json Failed")
	}

	//file.Read()
	//r := bufio.NewReader(file)
	//r.Rea
	//使用哪个IO操作想清楚

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Congfiguration)
	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error": err,
		}).Info("Decode config.json failed")
	}
}
