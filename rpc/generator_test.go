package rpc

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandStruct(t *testing.T) {
	filename := "/Users/ganxiangle/Desktop/workspace/scheduler/core/run.go"
	outfileNameServer := strings.TrimSuffix(filename, ".go") + ".server.go"
	outfileNameClient := strings.TrimSuffix(filename, ".go") + ".client.go"

	dirPath := filepath.Dir(filename)
	_, packageName := filepath.Split(dirPath)

	outFileForServer, err := os.Create(outfileNameServer)
	if err != nil {
		t.Fatal(err)
	}

	outFileForClient, err := os.Create(outfileNameClient)
	if err != nil {
		t.Fatal(err)
	}

	srcFile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(srcFile)
	if err != nil {
		t.Fatal(err)
	}
	srcFile.Close()
	srcFile, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}

	err = Generate(packageName, "", data, srcFile, outFileForServer, outFileForClient)
	if err != nil {
		t.Fatal(err)
	}
	srcFile.Close()
	outFileForServer.Close()
	outFileForClient.Close()

	err = exec.Command("go", "fmt", filename).Run()
	if err != nil {
		t.Fatal(err)
	}

	err = exec.Command("go", "fmt", outfileNameServer).Run()
	if err != nil {
		t.Fatal(err)
	}

	err = exec.Command("go", "fmt", outfileNameClient).Run()
	if err != nil {
		t.Fatal(err)
	}

}
