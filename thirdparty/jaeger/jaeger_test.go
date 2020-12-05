package jaeger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	tracing_log "github.com/opentracing/opentracing-go/log"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"log"
	"sacf-rpcx/common/jaeger/rpcxtracing"
	"sync"
	"testing"
	"time"
)

var serverAddr = "localhost:8972"
var jaegerAgentAddr = ""
var collectorEndpointAddr = "http://127.0.0.1:14268/api/traces"

type Request struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Request, reply *Reply) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Mul")
	span.LogFields(tracing_log.String("A", fmt.Sprintf("%d", args.A)))
	span.LogFields(tracing_log.String("B", fmt.Sprintf("%d", args.B)))
	time.Sleep(1 * time.Millisecond)
	defer span.Finish()

	reply.C = args.A * args.B
	return nil
}

func doClientRequest() {
	d := client.NewPeer2PeerDiscovery("tcp@"+serverAddr, "")
	option := client.DefaultOption
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	p := &client.OpenTracingPlugin{}
	pc := client.NewPluginContainer()
	pc.Add(p)
	xclient.SetPlugins(pc)

	args := &Request{A: 10, B: 20}
	reply := &Reply{}
	ctx := context.WithValue(
		context.Background(),
		share.ReqMetaDataKey,
		map[string]string{"msg": "from client"},
	)

	ctx = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
	err := xclient.Call(ctx, "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
}

func startServer() {
	s := server.NewServer()
	s.Plugins.Add(rpcxtracing.OpenTracingPlugin{})
	_ = s.RegisterName("Arith", new(Arith), "")
	_ = s.Serve("tcp", serverAddr)
}

func TestJaegerUsage(t *testing.T) {
	BufferFlushInterval = 1 * time.Second
	LogSpan = true

	InitJaeger("demo",
		"const",
		1,
		jaegerAgentAddr,
		collectorEndpointAddr,
		[]opentracing.Tag{{"test", "1"}},
	)

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
	time.Sleep(BufferFlushInterval * 2)
}
