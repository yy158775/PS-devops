package main

import (
	"chat-redis/redisservice"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
)

func main() {
	var wg sync.WaitGroup //不初始化的时候
	for i := 1;i <= 300;i ++ {
		wg.Add(1)
		go func() {
			clientconn,_ := grpc.Dial(":6380",grpc.WithBlock(),grpc.WithInsecure())
			defer clientconn.Close()
			defer wg.Done()
			client := redisservice.NewRedisServiceClient(clientconn)
			user,err := client.GetUserByName(context.Background(),&redisservice.UserName{UserName: "yy"})
			if err != nil {
				fmt.Println("error:",err)
				return
			}
			fmt.Println(user.UserName)
			//client.GetUserByName(context.Background(),&redisservice.UserName{})
			//client.GetUserByName(context.Background(),&redisservice.UserName{})
			//client.GetUserByName(context.Background(),&redisservice.UserName{})
		}()
	}
	wg.Wait()
}


