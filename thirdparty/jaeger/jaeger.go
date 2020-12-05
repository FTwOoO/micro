package jaeger

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

var (
	LogSpan             = true
	BufferFlushInterval = 10 * time.Second
)

func CreateJaegerTracer(cf *config.Configuration) (opentracing.Tracer, io.Closer, error) {
	if cf.Sampler == nil {
		cf.Sampler = &config.SamplerConfig{}
	}

	if cf.Reporter == nil {
		cf.Reporter = &config.ReporterConfig{}
	}

	tracer, closer, err := cf.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		err = fmt.Errorf("cannot initialize jaeger tracer: %w", err)
		return nil, nil, err
	}

	return tracer, closer, nil
}

func InitJaeger(serviceName string, sampleType string, sampleParam float64, agentAddr string, collectorEndpointAddr string, tags []opentracing.Tag) (io.Closer, error) {
	jaegerConfig := &config.Configuration{
		ServiceName: serviceName,
		Tags:        tags,
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: sampleParam,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort:  agentAddr,
			BufferFlushInterval: BufferFlushInterval,
			LogSpans:            LogSpan,
			CollectorEndpoint:   collectorEndpointAddr,
		},
	}

	tracer, closer, err := CreateJaegerTracer(jaegerConfig)
	if err != nil {
		return nil, err
	} else {
		opentracing.SetGlobalTracer(tracer)
		return closer, nil
	}
}
