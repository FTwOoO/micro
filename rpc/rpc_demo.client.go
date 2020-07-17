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

type DemoSeviceClient struct {
	*http.Client
	endPoint string
}

func NewDemoSeviceClient(c *http.Client, endpoint string) *DemoSeviceClient {
	return &DemoSeviceClient{
		Client:   c,
		endPoint: endpoint,
	}
}

func (c *DemoSeviceClient) doPostJsonAndUnpackRespJson(URL string, header http.Header, params interface{}, respObjPointer interface{}) (err error) {
	data, err := c.doPostJSON(URL, header, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, respObjPointer)
	return
}

func (c *DemoSeviceClient) doPostJSON(URL string, header http.Header, params interface{}) ([]byte, error) {
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

func (c *DemoSeviceClient) readResponse(resp *http.Response) ([]byte, error) {
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

func (this *DemoSeviceClient) SetEndPoint(end string) {
	this.endPoint = end
}

func (this *DemoSeviceClient) Hello(ctx context.Context, req *HelloRequest) (resp *HelloResponse, err error) {
	url := fmt.Sprintf("http://%s%s", this.endPoint, "/demoSevice/hello")
	resp = new(HelloResponse)
	httpResp := &demoSeviceHTTPResp{
		Data: resp,
	}
	err = this.doPostJsonAndUnpackRespJson(url, nil, req, httpResp)
	return
}
