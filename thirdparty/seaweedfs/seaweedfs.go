package seaweedfs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/linxGnu/goseaweedfs"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var sw *goseaweedfs.Seaweed

var FileNotFound = errors.New("file not found")
var filers []string

func Init(masterURL string, filers []string) (err error) {
	filers = filers
	sw, err = goseaweedfs.NewSeaweed(masterURL, filers, 8096, &http.Client{Timeout: 5 * time.Minute})
	if err != nil {
		return
	}
	return nil
}

func randomFiler() *goseaweedfs.Filer {
	filers := sw.Filers()
	randomIndex := rand.Intn(len(filers))
	filer := filers[randomIndex]
	return filer
}

func PutWithReader(ctx context.Context, content io.Reader, fileSize int, path string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "seaweedfsPut")
	defer span.Finish()

	result, err := randomFiler().Upload(content, int64(fileSize), path, "", "")
	if err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf("goseaweedfs upload fail:%s", result.Error)
	}
	return nil
}

func Put(ctx context.Context, content []byte, path string) (err error) {
	buf := bytes.NewBuffer(content)
	return PutWithReader(ctx, buf, len(content), path)
}

//TODO: 还应该返回content-type以便上层使用，但是因为第三方库没有暴露出这些信息，
//必要的时候需要修改第三方库实现
func Get(ctx context.Context, path string) (content []byte, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seaweedfsGet")
	defer span.Finish()
	span.LogFields(log.String("path", path))

	var code int
	content, code, err = randomFiler().Get(path, nil, nil)
	if code == http.StatusNotFound {
		return nil, FileNotFound
	}
	return
}

//TODO: 还应该返回content-type/size以便上层使用，但是因为第三方库没有暴露出这些信息，暂时不能返回
func Download(ctx context.Context, path string, cb func(io.Reader) error) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seaweedfsDownload")
	defer span.Finish()
	span.LogFields(log.String("path", path))
	return randomFiler().Download(path, nil, cb)
}

func Delete(ctx context.Context, path string) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seaweedfsDelete")
	defer span.Finish()
	span.LogFields(log.String("path", path))
	return randomFiler().Delete(path, nil)
}
