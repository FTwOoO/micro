package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

type demoImp struct {
}

func (this *demoImp) Hello(_ context.Context, req *HelloRequest) (resp *HelloResponse, err error) {
	fmt.Printf("req:%d]n", req.Send)
	return &HelloResponse{Ok: true}, nil
}

func (this *demoImp) HelloNoReponse(_ context.Context, req *HelloRequest) (err error) {
	fmt.Printf("req:%d]n", req.Send)
	return nil
}

func TestRegisterDemoSeviceForHTTP(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	RegisterDemoSeviceForHTTP(&demoImp{})

	go http.Serve(l, http.DefaultServeMux)
	time.Sleep(3 * time.Second)

	adr := fmt.Sprintf("127.0.0.1:%d", l.Addr().(*net.TCPAddr).Port)
	demoClient := NewDemoSeviceClient(http.DefaultClient, adr)
	resp, err := demoClient.Hello(context.TODO(), &HelloRequest{
		Send: 11,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Ok != true {
		t.Fatal("call fail")
	}

	err = demoClient.HelloNoReponse(context.TODO(), &HelloRequest{
		Send: 11,
	})
	if err != nil {
		t.Fatal(err)
	}

}
