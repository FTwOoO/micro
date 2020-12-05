package jaeger

import (
	"fmt"
	pb "github.com/FTwOoO/micro/thirdparty/jaeger/helloworld"
	"github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewClient() (pb.GreeterClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to service: %v", err)
	}
	return pb.NewGreeterClient(conn), nil
}
