package httpserver

import (
	"context"
	"gitlab.livedev.shika2019.com/go/common/cfg"
	"testing"
)

func TestHTTPServer_WaitHTTPServiceUp(t *testing.T) {
	server := NewHTTPServer(context.Background(), cfg.TestingEnv, &cfg.ConfigurationImp{
		HTTP: cfg.HTTPConfig{
			Addr:"0:0",

		},
	})
	server.Start()
	addr := server.WaitHTTPServiceUp()

	t.Log(addr)
}