package main

import (
	"bufio"
	"context"
	"grpc-demo/logger"
	"grpc-demo/message"
	"os"
)

func AfterLogin(c message.UserServiceClient, ch chan<- int) {
	logger.Info("\n----------------login succeed!----------------\n")

	for {
		logger.Info("\t\tyou can send group message or type exit:\n")
		//logger.Info("\t\tyou can exit by type exit\n")
		inputReader := bufio.NewReader(os.Stdin)
		content, err := inputReader.ReadString('\n')
		if content == "exit" {
			ch <- 1
			logger.Info("Exit")
			return
		}
		chat := message.ChatMessage{
			UserName: loginuser.UserName,
			Data:     content,
		}
		_, err = c.SendMessage(context.Background(), &chat)
		if err != nil {
			logger.Warn("SendMessage Failed")
			continue
		}
	}
}
