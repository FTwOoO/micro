package server

import (
	"os"
	"testing"
)

func TestIP(t *testing.T) {
	os.Setenv("INTRANET_IP", "127.0.0.2")
	ip := InternalIP()
	if ip != "127.0.0.2" {
		t.FailNow()
	}
	t.Logf("get ip:%v", ip)
}
