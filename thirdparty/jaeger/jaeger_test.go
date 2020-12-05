package jaeger

import (
	"context"
	pb "github.com/FTwOoO/micro/thirdparty/jaeger/helloworld"
	"github.com/opentracing/opentracing-go"
	"sync"
	"testing"
	"time"
)

var jaegerAgentAddr = ""
var collectorEndpointAddr = "http://127.0.0.1:14268/api/traces"
var collectorEndpointAddrForTest = "http://tracing-analysis-dc-sz.aliyuncs.com/adapt_dy016b2f2s@52c365727fcf1d5_dy016b2f2s@53df7ad2afe8301/api/traces"

func runClientAndServer(doClientRequest func(), startServer func()) {
	BufferFlushInterval = 1 * time.Second

	waitClient := &sync.WaitGroup{}
	waitClient.Add(1)

	waitServer := &sync.WaitGroup{}
	waitServer.Add(1)

	go func() {
		waitServer.Wait()
		doClientRequest()
		waitClient.Done()
	}()

	go startServer()
	waitServer.Done()
	waitClient.Wait()
	time.Sleep(1 * time.Second * 2)
}

func TestJaegerForGrpc(t *testing.T) {
	LogSpan = true

	InitJaeger("demo",
		"const",
		1,
		jaegerAgentAddr,
		collectorEndpointAddrForTest,
		[]opentracing.Tag{{"test", "1"}},
	)
	runClientAndServer(func() {
		c, _ := NewClient()
		c.SayHello(context.Background(), &pb.HelloRequest{Name: "123"})
	}, StartGrpcServer)

	runClientAndServer(runHTTPClient, runHTTPServer)
}
