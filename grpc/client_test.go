// @Author : Lik
// @Time   : 2021/2/1
package main

import (
	"GCSprout/grpc/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const (
	address = "localhost:16688"
)

func userSaveAll(client user.UserClient) {
	userList := make([]user.UserEntity, 2)
	userList[0] = user.UserEntity{
		Name: "库克",
		Age:  12,
	}
	userList[1] = user.UserEntity{
		Name: "雷军",
		Age:  14,
	}

	stream, err := client.UserSaveAll(context.Background())
	if err != nil {

	}

	finishChannel := make(chan bool)

	go func() {
		for true {
			res, err := stream.Recv()
			if err == io.EOF {
				finishChannel <- true
				break
			}
			if err != nil {

			}
			fmt.Println(res.User)
		}
	}()

	for _, u := range userList {
		stream.Send(&user.UserRequest{User: &u})
		if err != nil {

		}
	}
	stream.CloseSend()
	<-finishChannel
}

func userPhoto(client user.UserClient) {
	file, err := os.Open("12.jpg")
	if err != nil {
	}
	defer file.Close()
	md := metadata.New(map[string]string{"name": "张三"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.UserPhoto(ctx)
	if err != nil {

	}

	for {
		chunk := make([]byte, 128*1024)
		chunkSize, err := file.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {

		}

		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		stream.Send(&user.UserPhotoRequest{Data: chunk})
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {

	}
	fmt.Println(res.IsOk)

}

func userIndex(ctx context.Context, client user.UserClient) {
	//发起请求
	var req user.UserIndexRequest
	req.PageSize = 10
	req.Page = 1

	userIndexResp, err := client.UserIndex(ctx, &req)
	if err != nil {
		// 获取错误状态
		stat, ok := status.FromError(err)
		if ok {
			// 判断是否为调用超时
			if stat.Code() == codes.DeadlineExceeded {
				log.Fatalln("userIndex request timeout!")
			}
		}
		log.Fatalf("user index could not greet: %v\n", err)
	}

	if 0 == userIndexResp.Err {
		log.Printf("user index success: %s\n", userIndexResp.Msg)
		log.Printf("resp data is: %v", userIndexResp.Data)
	}

}

func userList(client user.UserClient) {

	stream, err := client.UserList(context.Background(), &user.UserListRequest{})
	if err != nil {

	}
	for true {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("receive userlist failed: %v", err)
		}
		fmt.Println(res.User)
	}

}

func TestTLSClient(t *testing.T) {

	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatalf("failed to tls: %v\n", err)
	}

	// 构建Token
	token := Token{
		AppId:     "grpc_token",
		AppSecret: "123456",
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&token),
	}

	//建立连接
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}

	defer conn.Close()

	userClient := user.NewUserClient(conn)

	// 设置请求超时
	// gRPC默认的请求的超时时间是很长的，当你没有设置请求超时时间时，
	// 所有在运行的请求都占用大量资源且可能运行很长的时间，导致服务
	// 资源损耗过高，使得后来的请求响应过慢，甚至会引起整个进程崩溃。
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//5:运行出现错误解决办法
	//
	//出现错误：
	//
	//2020/12/08 16:41:22 could not greet: rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0"
	//Process exiting with code: 1 signal: false
	//
	//
	//解决办法：GODEBUG=x509ignoreCN=0 go run client.go
	// 原因
	// go.15 以上已经启用了 CommonName了，因此推荐使用 SAN 证书.如果采用上述的方式生成的证书需要添加GODEBUG=x509ignoreCN=0 这样的环境变量才能正常运行

	userIndex(ctx, userClient)
	fmt.Println("--------------------")
	userList(userClient)
	fmt.Println("--------------------")
	userPhoto(userClient)
	fmt.Println("--------------------")
	userSaveAll(userClient)
}

func TestClient(t *testing.T) {

}
