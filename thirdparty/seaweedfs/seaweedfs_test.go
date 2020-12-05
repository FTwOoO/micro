package seaweedfs

import (
	"bytes"
	"context"
	"testing"
)

func init() {
	err := Init("http://127.0.0.1:9333", []string{"http://127.0.0.1:18888"})
	if err != nil {
		panic(err)
	}
}

func TestPut(t *testing.T) {
	path := "test/hello.txt"
	content := []byte("hello world")
	err := Put(context.TODO(), content, path)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("upload file to %s", path)

	contentDownload, err := Get(context.TODO(), path)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(contentDownload, content) != 0 {
		t.Fatal("content download != content upload")
	}
}

func TestDelete(t *testing.T) {
	path := "test/hello.txt"
	content := []byte("hello world")
	err := Put(context.TODO(), content, path)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("upload file to %s", path)

	err = Delete(context.TODO(), path)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Get(context.TODO(), path)
	if err != FileNotFound {
		t.Fatal(err)
	}

}
