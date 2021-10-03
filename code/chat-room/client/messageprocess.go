package main

import (
	"encoding/json"
	"fmt"
	"grpc-demo/logger"
	"grpc-demo/model"
	"net"
)

func Response(conn net.Conn, ch <-chan int) {
	buf := make([]byte, 1000)
	defer conn.Close()
	for {
		//not block
		select {
		case <-ch:
			logger.Info("Response Return")
			return
		default:
		}

		//通道阻塞和非阻塞的研究
		//学过的东西不用肯定忘记得很快
		//logger.Info("begin to receive")
		n, err := conn.Read(buf)
		//logger.Info("User receive:",buf)
		if err != nil {
			logger.Info("Response Read failed")
			return
			//连接断开
		}
		// 弄个通道，一旦我那边要结束
		// 这边立马结束
		info := model.MessageResp{}
		json.Unmarshal(buf[:n], &info)
		fmt.Println("From: ", info.UserName)
		fmt.Println("Data: ", info.Data)
	}
}
