package rpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type demoSeviceHTTPResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type demoSeviceErrorI interface {
	Error() string
	GetMsg() string
	GetCode() int
}

func RegisterDemoSeviceForHTTP(service DemoSevice) {

	http.DefaultServeMux.HandleFunc("/demoSevice/hello", func(w http.ResponseWriter, r *http.Request) {
		var httpStatusCode = http.StatusOK
		var code int
		var msg string
		var err error

		req := new(HelloRequest)
		body, _ := ioutil.ReadAll(r.Body)

		if err = json.Unmarshal(body, req); err != nil {
			httpStatusCode = http.StatusBadRequest
			code = http.StatusBadRequest
			msg = err.Error()
		}

		var reqI interface{} = req
		if v, ok := reqI.(interface{ Validate() error }); ok {
			err := v.Validate()
			if err != nil {
				httpStatusCode = http.StatusBadRequest
				code = http.StatusBadRequest
				msg = err.Error()
			}
		}

		var resp *HelloResponse
		if code == 0 {
			inCtx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
			resp, err = service.Hello(inCtx, req)

			if err != nil {
				if v, ok := err.(demoSeviceErrorI); ok {
					code = v.GetCode()
					msg = v.GetMsg()
				} else {
					code = http.StatusInternalServerError
					msg = err.Error()
				}
			} else {
				code = 0
				msg = ""
			}
		}

		respD, _ := json.Marshal(demoSeviceHTTPResp{
			Code: code,
			Msg:  msg,
			Data: resp,
		})

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(httpStatusCode)
		w.Write([]byte(respD))
	})

}
