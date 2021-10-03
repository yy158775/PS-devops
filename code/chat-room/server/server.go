package main

import (
	//"bufio"
	"google.golang.org/grpc"
	"grpc-demo/message"
	"grpc-demo/server/config"
	"grpc-demo/server/logrus"
	"grpc-demo/server/model"
	"grpc-demo/server/utils"
	"log"
	"net"
)


func messageprocess(lis net.Listener) {
	for {  //接受别人的链接请求
		conn,err := lis.Accept()
		if err != nil {
			continue
		}
		buffer := make([]byte,500)
		//token,_ := bufio.NewReader(conn).ReadString(0)
		//var buffer bytes.Buffer
		n,_ := conn.Read(buffer)
		//io.Copy(&buffer,conn) //eof是不是只能关了时候，才会出现
		//logrus.Log.Info("len: ",buffer.Len())

		logrus.Log.Info(string(buffer[:n]))
		logrus.Log.Info(len(string(buffer[:n])))
		//userinfo ,err := utils.ValidateToken(token)
		//userinfo ,err := utils.ValidateToken(buffer.String())
		userinfo ,err := utils.ValidateToken(string(buffer[0:n]))
		if err != nil {  //验证失败
			logrus.Log.WithFields(map[string]interface{}{
				"Error" : err,
			}).Warn("Jwt validation failed")
			continue
		} else {
			logrus.Log.WithFields(map[string]interface{}{
				"User" : userinfo.UserName,
			}).Info("Login Succeed")
		}
		model.RWmutex.Lock()
		model.ClientConnsMap[userinfo.UserName] = model.ConnInfo{
			Conn:   conn,
			Status: true,
		}
		model.RWmutex.Unlock()
		logrus.Log.WithFields(map[string]interface{}{
			"User":userinfo.UserName,
		}).Info("Login Succeed")
	}
}

func main() {
	server := grpc.NewServer()
	message.RegisterUserServiceServer(server,new(model.ChatServer))
	rpcconn,err := net.Listen("tcp",config.Congfiguration.ServerInfo.RpcHost) //这个负责rpc远程函数调用三个功能，
	//登录，注册，发消息

	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error": err,
		}).Fatal("Listening ",config.Congfiguration.ServerInfo.RpcHost)
	} else {
		logrus.Log.Info("Listening Rpc",config.Congfiguration.ServerInfo.RpcHost)
	}

	lisconn,err := net.Listen("tcp",config.Congfiguration.ServerInfo.RespHost) //这个负责接受客户端的请求，通过这个链接给客户端返回消息

	if err != nil {
		logrus.Log.WithFields(map[string]interface{}{
			"Error": err,
		}).Fatal("Listening Resp",config.Congfiguration.ServerInfo.RespHost)
	} else {
		logrus.Log.Info("Listening ",config.Congfiguration.ServerInfo.RespHost)
	}
	go messageprocess(lisconn)
	//server.Serve(rpcconn)
	if err := server.Serve(rpcconn); err != nil {
		log.Fatalf("Failed to serve port %v : %v", rpcconn,err)
	}
}

