package server

import (
	"context"
	"github.com/rexue2019/micro/cfg"
	"testing"
)

func TestHTTPServer_WaitHTTPServiceUp(t *testing.T) {
	server := NewHTTPServer(context.Background(), cfg.TestingEnv, &cfg.ConfigurationImp{
		HTTP: cfg.HTTPConfig{
			Addr: "0:0",
		},
	})
	server.Start()
	addr := server.WaitHTTPServiceUp()

	t.Log(addr)
}
