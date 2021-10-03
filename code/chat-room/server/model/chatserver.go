package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"google.golang.org/grpc"
	"grpc-demo/common"
	"grpc-demo/logger"
	"grpc-demo/message"
	"grpc-demo/model"
	redisservice "grpc-demo/redisclient"
	"grpc-demo/server/config"
	"grpc-demo/server/logrus"
	"grpc-demo/server/utils"
	"net"
	"sync"
)

type ConnInfo struct {
	Conn   net.Conn
	Status bool
}

var RWmutex sync.RWMutex
var ClientConnsMap = make(map[string]ConnInfo)

type ChatServer struct {
	//map[string]
}

func (c *ChatServer) Register(ctx context.Context, req *message.RegisterRequest) (res *message.RegisterResponse, err error) {
	res = &message.RegisterResponse{} //才发现这个问题
	if req.ConfirmPassWord != req.PassWord {
		res.Code = common.PasswordNotMatch
		res.Message = common.ERROR_PASSWORD_DOES_NOT_MATCH.Error()
		return res, common.ERROR_PASSWORD_DOES_NOT_MATCH //这一点
	}

	//链接redis
	conn, err := grpc.Dial(config.Congfiguration.RedisInfo.RedisServerHost, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error":  err,
			"Target": config.Congfiguration.RedisInfo.RedisServerHost, //ip 地址
		}).Fatal("Redis Connection Error")
		return res, err
	} else {
		logrus.Log.Info("Redis Connection Succeed")
	}

	client := redisservice.NewRedisServiceClient(conn)

	_, err = client.GetUserByName(context.Background(), &redisservice.UserName{
		UserName: req.UserName,
	})
	// 保证用户名不重复 这一点我不清楚
	// 没有 和 出错是否是一个概念
	if err == nil {
		logrus.Log.WithFields(map[string]interface{}{
			"UserName": req.UserName,
		}).Info("UserHasExisted")
		res.Code = common.UserHasExited //已存在
		res.Message = common.ERROR_USER_ALREADY_EXISTS.Error()
		return res, common.ERROR_USER_ALREADY_EXISTS //错误处理怎么办
	}
	h := sha256.New()
	h.Write([]byte(req.PassWord))
	newuser := redisservice.NewUser{
		UserName: req.UserName,
		Password: hex.EncodeToString(h.Sum(nil)), //加密
	}

	_, err = client.InsertUser(ctx, &newuser) //是不是要加返回信息
	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error":  err,
			"Target": config.Congfiguration.RedisInfo.RedisServerHost, //ip 地址
		}).Warn("Redis Insert Error")
		return res, err
	}

	res.Code = common.LoginSucceed
	res.Message = common.Register_Success //返回信息怎么写
	return res, nil
}

func (c *ChatServer) Login(ctx context.Context, req *message.LoginRequest) (res *message.LoginResponse, err error) {

	res = &message.LoginResponse{} //才发现这个问题
	conn, err := grpc.Dial(config.Congfiguration.RedisInfo.RedisServerHost, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error":  err,
			"Target": config.Congfiguration.RedisInfo.RedisServerHost, //ip 地址
		}).Fatal("Redis Connection Error")
		return res, err
	}

	client := redisservice.NewRedisServiceClient(conn)
	username := redisservice.UserName{UserName: req.UserName}

	userinfo, err := client.GetUserByName(context.Background(), &username)

	//用户不存在？ 还是没连上 这种错误怎么处理恩？ 查不到会怎么办，返回值如何处理

	if err != nil {
		res.Code = common.UserNotExist
		res.Message = common.ERROR_USER_DOES_NOT_EXIST.Error()
		return res, common.ERROR_USER_DOES_NOT_EXIST
	}

	h := sha256.New()
	h.Write([]byte(req.PassWord))

	if userinfo.Password == hex.EncodeToString(h.Sum(nil)) { //登陆成功
		res.Code = common.LoginSucceed
		//ClientConnsMap[userinfo.UserName] = ConnInfo{Status: true,Conn: nil}
		tokenString := utils.GenerateToken(&utils.User{UserName: req.UserName}, 0)
		res.Message = tokenString //生成一个JWT
		return res, nil
	} else {
		res.Code = common.PasswordNotMatch
		res.Message = common.ERROR_PASSWORD_DOES_NOT_MATCH.Error() //这个状态码也设置一下吧
		return res, common.ERROR_PASSWORD_DOES_NOT_MATCH
	}
}

func (c *ChatServer) SendMessage(ctx context.Context, chatmessage *message.ChatMessage) (*message.Empty, error) {
	//得加锁吧
	RWmutex.RLock()
	for key, value := range ClientConnsMap {
		if !value.Status {
			delete(ClientConnsMap, key)
		} else {
			sendmessage := model.MessageResp{
				UserName: chatmessage.GetUserName(),
				Data:     chatmessage.GetData(),
			}
			buffer, err := json.Marshal(sendmessage)
			if err != nil {
				continue
			}
			logger.Info("forward message:", string(buffer))
			_, err = value.Conn.Write(buffer)
			if err != nil {
				delete(ClientConnsMap, key)
				continue
			}
		}
	}
	RWmutex.RUnlock()
	return &message.Empty{}, nil
}
