package examples

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

func NewJaegerTracer(service string) opentracing.Tracer {
	sender := transport.NewHTTPTransport(
		// 设置链路跟踪的网关（不同region对应不同的值，从http://tracing.console.aliyun.com/ 的配置查看中获取）
		"http://tracing-analysis-dc-sz.aliyuncs.com/adapt_dy016b2f2s@52c365727fcf1d5_dy016b2f2s@53df7ad2afe8301/api/traces",
	)
	tracer, _ := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),
	)
	return tracer
}
