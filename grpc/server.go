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
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

const (
	port = "localhost:16688"
)

type UserService struct {
}

var userDataList = []user.UserEntity{
	{Name: "张三", Age: 13},
	{Name: "王五", Age: 12},
	{Name: "杰克", Age: 18},
	{Name: "坤坤", Age: 15},
}

func handle(ctx context.Context, req *user.UserIndexRequest, data chan<- *user.UserIndexResponse) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		// 超时后退出该携程
		runtime.Goexit()
	case <-time.After(4 * time.Second): // 模拟耗时操作
		res := user.UserIndexResponse{
			Err: 0,
			Msg: "success",
			Data: []*user.UserEntity{
				{Name: "Jack", Age: 12},
				{Name: "Tom", Age: 15},
			},
		}
		// //修改数据库前进行超时判断
		// if ctx.Err() == context.Canceled{
		// 	...
		// 	//如果已经超时，则退出
		// }
		data <- &res
	}

}

func (s *UserService) UserIndex(ctx context.Context, req *user.UserIndexRequest) (*user.UserIndexResponse, error) {
	log.Printf("receive user index request: page %d page_size %d", req.Page, req.PageSize)

	data := make(chan *user.UserIndexResponse, 1)
	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
}

func (s *UserService) UserView(ctx context.Context, in *user.UserViewRequest) (*user.UserViewResponse, error) {
	log.Printf("receive user view request: uid %d", in.Uid)

	return &user.UserViewResponse{
		Err:  0,
		Msg:  "success",
		Data: &user.UserEntity{Name: "james", Age: 28},
	}, nil
}

func (s *UserService) UserPost(ctx context.Context, in *user.UserPostRequest) (*user.UserPostResponse, error) {
	log.Printf("receive user post request: name %s password %s age %d", in.Name, in.Password, in.Age)

	return &user.UserPostResponse{
		Err: 0,
		Msg: "success",
	}, nil
}

func (s *UserService) UserDelete(ctx context.Context, in *user.UserDeleteRequest) (*user.UserDeleteResponse, error) {
	log.Printf("receive user delete request: uid %d", in.Uid)

	return &user.UserDeleteResponse{
		Err: 0,
		Msg: "success",
	}, nil
}

func (s *UserService) UserList(request *user.UserListRequest, stream user.User_UserListServer) error {
	log.Println("receive user list request")
	for _, e := range userDataList {
		stream.Send(&user.UserResponse{
			User: &e,
		})
	}
	return nil
}

func (s *UserService) UserPhoto(stream user.User_UserPhotoServer) error {
	log.Println("receive user photo request")

	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("User: %s\n", md["name"][0])
	}

	var img []byte
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File Size: %d\n", len(img))
			return stream.SendAndClose(&user.UserPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("File reveived: %d\n", len(res.Data))
		img = append(img, res.Data...)
	}

	return nil
}

func (s *UserService) UserSaveAll(stream user.User_UserSaveAllServer) error {
	log.Println("receive user saveAll request")
	for true {
		userReq, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		userDataList = append(userDataList, *userReq.User)
		_ = stream.Send(&user.UserResponse{User: userReq.User})
	}

	for _, user := range userDataList {
		fmt.Println(user)
	}
	return nil

}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 生成证书
	// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj /CN=localhost
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("failed to tls server: %v", err)
	}

	// 添加TLS认证和Token认证
	options := []grpc.ServerOption{
		grpc.Creds(creds),
		// 当存在多个拦截器时，只取第一个拦截器
		grpc.UnaryInterceptor(ServerTokenInterceptor),
		//grpc.UnaryInterceptor(),
	}

	//创建grpc服务容器
	grpServer := grpc.NewServer(options...)

	// 为user服务注册业务实现	 将user服务绑定到rpc服务容器上
	user.RegisterUserServer(grpServer, &UserService{})
	// 注册反射服务 这个服务器是Cli使用的 跟服务本身没有关系
	reflection.Register(grpServer)

	log.Println("服务正在启动...")

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	if err := grpServer.Serve(listen); err != nil {
		log.Fatalf("fialed to serve: %v", err)
	}

}

type Token struct {
	AppId     string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据（metadata）
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppId, "app_secret": t.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}

// Check 验证token
func CheckToken(ctx context.Context) error {
	//从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "获取Token失败")
	}
	var (
		appID     string
		appSecret string
	)
	if value, ok := md["app_id"]; ok {
		appID = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appID != "grpc_token" || appSecret != "123456" {
		return status.Errorf(codes.Unauthenticated, "Token无效: app_id=%s, app_secret=%s", appID, appSecret)
	}
	return nil
}

func ServerTokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 拦截普通方法
	err = CheckToken(ctx)
	if err != nil {
		return
	}

	return handler(ctx, req)
}

type AuthFunc func(ctx context.Context) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		newCtx.Done()
		//wrapped := grpc_middleware.WrapServerStream(stream)
		//wrapped.WrappedContext = newCtx
		//return handler(srv, wrapped)
		return handler(srv, nil)
	}
}
