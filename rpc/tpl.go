package rpc

const rpcServerCommonTpl = `
type {{ .serviceName }}HTTPResp struct { 
	Code int  ` + "`json:\"code\"`" + `
	Msg string  ` + "`json:\"msg\"`" + `
	Data interface{}` + "`json:\"data\"`}" + `

type {{ .serviceName }}ErrorI interface {
	Error() string
	GetMsg() string
	GetCode() int
}
`

const httpRouteHandlerTpl = `
 
	http.DefaultServeMux.HandleFunc("{{ .httpMethodName }}", func(w http.ResponseWriter, r *http.Request) {
		var httpStatusCode = http.StatusOK
		var code int
		var msg string
 		var err error

		req := new({{ .rpcRequestStructName }})
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

		var resp *{{ .rpcResponseStructName }}
		if code == 0 {
			inCtx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
			resp, err = service.{{ .rpcMethodName }}(inCtx, req)

			if err != nil {
				if v, ok := err.({{ .serviceName }}ErrorI); ok {
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

		respD, _ := json.Marshal({{ .serviceName }}HTTPResp{
			Code: code,
			Msg:  msg,
			Data: resp,
		})

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(httpStatusCode)
		w.Write([]byte(respD))
	})
 
`

const httpRouteHandlerForNoRespTpl = `
 
	http.DefaultServeMux.HandleFunc("{{ .httpMethodName }}", func(w http.ResponseWriter, r *http.Request) {
		var httpStatusCode = http.StatusOK
		var code int
		var msg string
 		var err error

		req := new({{ .rpcRequestStructName }})
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

		if code == 0 {
			inCtx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
			err = service.{{ .rpcMethodName }}(inCtx, req)

			if err != nil {
				if v, ok := err.({{ .serviceName }}ErrorI); ok {
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

		respD, _ := json.Marshal({{ .serviceName }}HTTPResp{
			Code: code,
			Msg:  msg,
			Data: nil,
		})

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(httpStatusCode)
		w.Write([]byte(respD))
	})
 
`

const RcpClientCommonTpl = `
package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type {{ .structName }} struct {
	*http.Client
	endPoint string
}

func New{{ .structName }}(c *http.Client, endpoint string) *{{ .structName }} {
	return &{{ .structName }}{
		Client: c,
		endPoint: endpoint,
	}
}

func (c *{{ .structName }}) doPostJsonAndUnpackRespJson(URL string, header http.Header, params interface{}, respObjPointer interface{}) (err error) {
	data, err := c.doPostJSON(URL, header, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, respObjPointer)
	return
}

func (c *{{ .structName }}) doPostJSON(URL string, header http.Header, params interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var content string

	switch params.(type) {
	case []byte:
		content = string(params.([]byte))
	case string:
		content = params.(string)
	default:

		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(params)
		if err != nil {
			return nil, err
		}

		content = buf.String()
	}

	body := strings.NewReader(content)
	req, err = http.NewRequest("POST", URL, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ct, _ := c.readResponse(resp)
		err = fmt.Errorf("Fail to request(%s): [%d] %s", URL, resp.StatusCode, string(ct))
		return nil, err
	}

	return c.readResponse(resp)
}

func  (c *{{ .structName }}) readResponse(resp *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = resp.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func (this *{{ .structName }}) SetEndPoint(end string) {
	this.endPoint = end
}`

const RpcClientMethodTpl = `
func (this *{{ .structName }}) {{ .rpcMethodName }}(ctx context.Context, req *{{ .rpcRequestStructName }}) (resp *{{ .rpcResponseStructName }}, err error) {
	url := fmt.Sprintf("http://%s%s", this.endPoint, "{{ .httpMethodName }}")
	resp = new({{ .rpcResponseStructName }})
	httpResp := &{{ .serviceName }}HTTPResp{
		Data: resp,
	}
	err = this.doPostJsonAndUnpackRespJson(url, nil, req, httpResp)
	return
}`

const RpcClientMethodForNoRespTpl = `
func (this *{{ .structName }}) {{ .rpcMethodName }}(ctx context.Context, req *{{ .rpcRequestStructName }}) (err error) {
	url := fmt.Sprintf("http://%s%s", this.endPoint, "{{ .httpMethodName }}")
	_, err = this.doPostJSON(url, nil, req)
	return
}`
