package main

import (
	"chat-redis/config"
	"chat-redis/redisservice"
	"context"
	"github.com/garyburd/redigo/redis"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	redisservice.UnimplementedRedisServiceServer
}

func (s *Server) GetUserByName(ctx context.Context, req *redisservice.UserName) (userinfo *redisservice.UserInfo, err error) {
	userinfo = &redisservice.UserInfo{}
	conn := config.Pool.Get()
	defer conn.Close()
	infobytes,err := redis.Bytes(conn.Do("hget", "users", req.UserName))

	if err != nil {
		log.Println("hget failed")
		return userinfo,err
	}
	//redis如何存储数据
	err = userinfo.XXX_Unmarshal(infobytes)
	if err != nil {
		return userinfo,err
	}
	return &redisservice.UserInfo{UserName: userinfo.UserName,Password: userinfo.Password},nil
	//思考一下该怎么存储 json 还是 其他的
}

func (s *Server) InsertUser(ctx context.Context, req *redisservice.NewUser) (*redisservice.Empty, error) {
	conn := config.Pool.Get()
	defer conn.Close()
	newuser := redisservice.NewUser{UserName: req.UserName,Password: req.Password}

	res,err := newuser.XXX_Marshal(nil,false)
	if err != nil {
		log.Println("XXX_Marshal failed")
		return &redisservice.Empty{},err
	}

	_, err = conn.Do("hset", "users", req.UserName, res)

	if err != nil {
		log.Println("redis hset failed",err)
		return &redisservice.Empty{},err
	}
	return &redisservice.Empty{},nil
}

func main() {
	server := grpc.NewServer()
	conn,err := net.Listen("tcp",config.RedisInfo.ListenHost)  //监听在某个端口
	if err != nil {
		log.Fatal("Redis Listen Error")
	}

	redisservice.RegisterRedisServiceServer(server,new(Server))
	//这里通过这个函数将server的解码编码给实现了注册到这个服务器上
	//这是proto那个实现的go，应该也可以用json进行注册

	err = server.Serve(conn)
	if err != nil {
		log.Fatal("Redis Server Error")
	}
}