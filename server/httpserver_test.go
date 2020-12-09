package server

import (
	"context"
	"github.com/FTwOoO/micro/cfg"
	"net/http"
	"testing"
	"time"
)

func TestHTTPServer_WaitHTTPServiceUp(t *testing.T) {
	server := NewHTTPServer(context.Background(), cfg.TestingEnv, &cfg.ConfigurationImp{
		HTTP: cfg.HTTPConfig{
			Addr: "0:11999",
		},
	})
	server.AddMiddleware(AHASMiddleware("7bac187365984241afadce74133c9820", "go-test-demo"))
	server.Start()
	addr := server.WaitHTTPServiceUp()
	go func() {
		for i := 0; i < 100000; i++ {
			_, _ = http.Get("127.0.0.1:11999")
			time.Sleep(1 * time.Millisecond)
		}
	}()
	time.Sleep(100 * time.Second)

	t.Log(addr)
}
