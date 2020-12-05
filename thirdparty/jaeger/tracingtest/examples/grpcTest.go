package examples

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"time"
	pb "tracingtest/examples/helloworld"
)

const (
	addr           = "0.0.0.0:19090"
	grpcServerName = "grpcServer"
	grpcClientName = "grpcClient"
)

func GrpcTest() {
	tracer := NewJaegerTracer(grpcServerName)
	opentracing.SetGlobalTracer(tracer)
	go StartGrpcServer(tracer)
	time.Sleep(1 * time.Second)
	c, _ := NewClient(NewJaegerTracer(grpcClientName))
	ctx := context.Background()
	c.SayHello(ctx, &pb.HelloRequest{Name: "123"})
	fmt.Println("grpc done")
}
