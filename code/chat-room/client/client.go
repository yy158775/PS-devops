package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-demo/logger"
	"grpc-demo/message"
	"grpc-demo/server/config"
	"io"
	"net"
)

type LoginUser struct {
	UserName    string
	tokenString string
}

var loginuser LoginUser

func main() {
	var (
		key             int
		loop            = true
		userName        string
		password        string
		confirmpassword string
	)
	// Set up a connection to the server.
	conn, err := grpc.Dial(config.Congfiguration.ServerInfo.RpcHost, grpc.WithInsecure(), grpc.WithBlock())
	//block 千万小心
	if err != nil {
		logger.Error("Dial %v failed:%v", config.Congfiguration.ServerInfo.RpcHost, err)
		return
	} else {
		logger.Info("Rpc Service is ready")
	}
	defer conn.Close()
	c := message.NewUserServiceClient(conn)

	for loop {

		logger.Info("\n----------------Welcome to the chat room--------------\n")
		logger.Info("\t\tSelect the options：\n")
		logger.Info("\t\t\t 1、Sign in\n")
		logger.Info("\t\t\t 2、Sign up\n")
		logger.Info("\t\t\t 3、Exit the system\n")

		// get user input
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			logger.Info("sign In Please\r\n")
			logger.Notice("Username:\n")
			fmt.Scanf("%s\n", &userName)
			logger.Notice("Password:\n")
			fmt.Scanf("%s\n", &password)

			//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			//defer cancel()

			resp, err := c.Login(context.Background(), &message.LoginRequest{UserName: userName, PassWord: password})

			if err != nil {
				//错误处理
				logger.Error("Login failed: %v\r\n", err)

			} else if resp.Code == 200 {
				receiveconn, err := net.Dial("tcp", config.Congfiguration.ServerInfo.RespHost)
				//链接一个新的连接，用于接受消息 另外一个端口了
				if err != nil {
					logger.Warn("Connect %s failed: %v\r\n", config.Congfiguration.ServerInfo.RespHost, err)
					continue
				}
				loginuser.tokenString = resp.Message
				loginuser.UserName = userName
				io.WriteString(receiveconn, loginuser.tokenString)
				//这个方法不错

				ch := make(chan int, 2)
				go Response(receiveconn, ch) //负责接受消息
				AfterLogin(c, ch)            //负责发送是否结束的消息
			}

		case 2:
			logger.Info("Create account\n")
			logger.Notice("user name：\n")
			fmt.Scanf("%s\n", &userName)
			logger.Notice("password：\n")
			fmt.Scanf("%s\n", &password)
			logger.Notice("password confirm：\n")
			fmt.Scanf("%s\n", &confirmpassword)

			resp, err := c.Register(context.Background(), &message.RegisterRequest{UserName: userName,
				PassWord: password, ConfirmPassWord: confirmpassword})

			if err != nil {
				logger.Error("Create account failed: %v\n", err)
				continue
			}
			//resp 如何处理
			if resp.Code == 200 {
				//注册成功
				logger.Info("Create account succeed")
			}

		case 3:
			logger.Warn("Exit...\n")
			loop = false // this is equal to 'os.Exit(0)'
		default:
			logger.Error("Select is invalid!\n")
		}
	}
}
