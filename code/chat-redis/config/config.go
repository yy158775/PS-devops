package config

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)


type redisInfo struct {
	Host        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	ListenHost string
}

var RedisInfo = redisInfo{}

func init() {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		fmt.Printf("Open file error: %v\n", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&RedisInfo)
	if err != nil {
		fmt.Println("Json Decoder Error: ", err)
		return
	}
}


var Pool *redis.Pool

func initRedisPool(maxIdle, maxActive int, idleTimeout time.Duration, host string) {
	Pool = &redis.Pool{
		// 初始化链接数量
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host)
		},
	}
	conn := Pool.Get()
	defer conn.Close()
}

func init() {
	// 初始化 redis 连接池，全局唯一  redis怎么配置是个问题
	redisInfo := RedisInfo
	fmt.Println("redisInfo", redisInfo)
	initRedisPool(redisInfo.MaxIdle, redisInfo.MaxActive, time.Second*(redisInfo.IdleTimeout), redisInfo.Host)
}
