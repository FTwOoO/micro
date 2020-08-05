package server

import (
	"context"
	"github.com/rexue2019/micro/cfg"
	"github.com/rexue2019/util/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MicroContext struct {
	Context       context.Context
	ConfigPointer cfg.Configuration
	HTTPServer    *HTTPServer
	Logger        logging.Logger
}

func startHttpServer(ctx context.Context, configPointer cfg.Configuration) *HTTPServer {
	httpServer := NewHTTPServer(ctx, configPointer.GetEnv(), configPointer)
	go httpServer.Start()
	addr := httpServer.WaitHTTPServiceUp()
	configPointer.GetHttp().Addr = addr

	if configPointer.GetHttp().EnablePrometheusMetrics {
		logging.Log.Infow(logging.KeyScope, "microkit", logging.KeyEvent, "enablePrometheusMetrics")
		httpServer.SetupPrometheusMetrics()
	}

	if configPointer.GetHttp().EnablePprof {
		logging.Log.Infow(logging.KeyScope, "microkit", logging.KeyEvent, "enablePprof")
		httpServer.SetupPprof()
	}
	return httpServer
}

//运行micro，同时会解释出config
func Run(configPointer cfg.Configuration, handler func(context MicroContext)) {
	var err error
	ctx := ContextWithSignal()

	cfg.ParseConfigFileFromFlags(configPointer)
	if configPointer.GetName() == "" || configPointer.GetVersion() == "" {
		logging.Log.Fatalw(logging.KeyScope, "microkit", logging.KeyMsg, "invalid config: missing name/version")
	}

	logging.SetLogSetting(configPointer.GetLogLevel(), true)

	var httpServer *HTTPServer
	if configPointer.GetHttp() != nil && configPointer.GetHttp().IsValid() {
		httpServer = startHttpServer(ctx, configPointer)
	}

	runCtx := MicroContext{
		Context:       ctx,
		ConfigPointer: configPointer,
		HTTPServer:    httpServer,
		Logger: logging.Log.WithFields(map[string]interface{}{
			logging.KeyService: configPointer.GetName(),
		}),
	}

	handler(runCtx)

	select {
	case <-runCtx.Context.Done():
	}

	err = httpServer.Close()
	if err != nil {
		logging.Log.Warnw(logging.KeyScope, "microkit", logging.KeyMsg, "close http server fail")

	}

	//把调度转给其它goroutine，这样就可以退出自动任务了，最好的做法是等待所有的goroutine退出
	time.Sleep(5 * time.Second)
	logging.Log.Infow(logging.KeyScope, "microkit", logging.KeyEvent, "serverExit")
	return
}

func ContextWithSignal() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		<-quit
		cancel()
	}()
	return ctx
}
