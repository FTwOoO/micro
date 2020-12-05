package jaeger

import (
	"context"
	"fmt"
	pb "github.com/FTwOoO/micro/thirdparty/jaeger/helloworld"
	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	addr = "0.0.0.0:19090"
)

func StartGrpcServer() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	// Register reflection service on gRPC server.
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())))
	pb.RegisterGreeterServer(s, &grpcServer{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}

type grpcServer struct{}

func (s *grpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		mysqlSpan, _ := opentracing.StartSpanFromContext(ctx, "SQL FindUserTable")
		mysqlSpan.SetTag("db.statement", "select * from user ...")
		//do mysql operations
		time.Sleep(time.Millisecond * 100)
		mysqlSpan.Finish()
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
