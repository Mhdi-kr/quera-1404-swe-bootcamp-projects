package grpcserver

import (
	"net"

	"example.com/authorization/internal/service"
	postv1 "example.com/authorization/protos-gen/post/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ListenAndServe(addr string, postSrv service.PostService) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	postv1.RegisterPostServiceServer(grpcServer, &PostServiceServer{postSrv: postSrv})
	reflection.Register(grpcServer)

	return grpcServer.Serve(listener)
}
