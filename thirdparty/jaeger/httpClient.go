package jaeger

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

func runHTTPClient() {
	c := &http.Client{Transport: &nethttp.Transport{}}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:%s/", serverPort),
		nil,
	)
	if err != nil {
		return
	}

	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req)
	defer ht.Finish()

	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	fmt.Printf("Received result: %s\n", string(body))
}

func onError(span opentracing.Span, err error) {
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print(err)
}
