package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rexue2019/micro/cfg"
	"github.com/rexue2019/util/errorkit"
	"github.com/rexue2019/util/logging"
	"go.uber.org/atomic"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

type HTTPServer struct {
	server        *http.Server
	httpConfig    *cfg.HTTPConfig
	start         atomic.Bool
	listenAddr    string
	handleConnect http.HandlerFunc
}

func (this *HTTPServer) WaitHTTPServiceUp() (addr string) {

	tk := time.Tick(10 * time.Millisecond)

	for {

		select {
		case <-tk:
			startSuccess := this.start.Load()

			if startSuccess {
				return this.listenAddr
			}
		}
	}
}

func NewHTTPServer(ctx context.Context, env cfg.Environment, config cfg.Configuration) *HTTPServer {
	service := &HTTPServer{httpConfig: config.GetHttp()}
	return service
}

func (this *HTTPServer) SetupPrometheusMetrics() {
	http.DefaultServeMux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
}

func (this *HTTPServer) SetupPprof() {
	_ = pprof.Index
}

func (this *HTTPServer) SetupConnectHandleFunc(f http.HandlerFunc) {
	this.handleConnect = f
}

func (this *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "CONNECT" && this.handleConnect != nil {
		this.handleConnect(w, req)
	} else {
		http.DefaultServeMux.ServeHTTP(w, req)
	}
}

func (this *HTTPServer) Start() {
	go func() {
		logging.Log.Infow("event", "httpServerStart", "addr", this.httpConfig.Addr)
		this.server = &http.Server{Addr: this.httpConfig.Addr, Handler: this}

		ln, err := net.Listen("tcp", this.httpConfig.Addr)
		if err != nil {
			logging.Log.Fatalw("event", "net.Listen fail")
		}
		this.listenAddr = InternalIP()
		_, port, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			panic(err)
		}
		this.listenAddr = InternalIP() + ":" + port

		this.start.Store(true)
		err = this.server.Serve(ln)
		if err != nil {
			logging.Log.LogError(errorkit.WrapError(err).AddOp("http.Server.Serve"))
		}
	}()

}

func (this *HTTPServer) Close() error {
	return this.server.Close()
}
