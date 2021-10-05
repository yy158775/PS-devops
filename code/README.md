# 项目名称：chat-room聊天室

原理：





# chat-service:

提供微服务接口：


rpc Register(RegisterRequest)      returns (RegisterResponse);

rpc Login(LoginRequest)       		returns (LoginResponse);

rpc SendMessage(ChatMessage)  returns (Empty);

# chat-redis:

提供的微服务接口

负责查询后端redis数据库，给server提供查询功能。

services截图

![Image-text](https://github.com/yy158775/PS-devops/blob/master/docs/photo/services.png)


# k8s运行

redis-service-58d68bcdfb-hfdg8    为redis镜像



# 自动化部署采用drone cloud

配置文件.drone.yml

# Quickly Start

## clone the repository

```
git clone https://github.com/yy158775/PS-devops
```

## config



code/chat-room/client/config.json

code/chat-room/server/config.json

```json
{
  "ServerInfo": {
    "RpcHost": "localhost:8080",
    "RespHost": "localhost:8081"
  },
  "RedisInfo": {
    "RedisServerHost": "localhost:6380"
  }
}
```



code/chat-redis/config/config.json

```json
{
  "Host": "redis-service:6379",
  "MaxIdle":     16,
  "MaxActive":   0,
  "IdleTimeout": 300,
  "ListenHost": ":6380"
}
```

## start redis

在本地服务器中启动redis-server，监听端口为默认端口6379

## start chat-redis

```
cd code/chat-redis/server
go run server.go
```

## start server

```
cd code/chat-room/server
go run server.go
```

## start client

```
cd code/chat-room/client
go build
./client
```

## the start sequence 

​	redis,chat-redis,server三者的启动顺序并不重要，只要保证client在这三者启动后即可。

​	因为只有在client发出请求时，微服务彼此之间才会产生调用，之前并未产生任何通信。

# Client

功能有：登录   注册    发言



