package main

import (
	"fmt"
	"github.com/FTwOoO/micro/rpc"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	filePaths := os.Args[1:]

	for _, filename := range filePaths {
		outfileNameServer := strings.TrimSuffix(filename, ".go") + ".server.go"
		outfileNameClient := strings.TrimSuffix(filename, ".go") + ".client.go"
		filename, _ = filepath.Abs(filename)
		dirPath := filepath.Dir(filename)
		_, packageName := filepath.Split(dirPath)
		fmt.Println("use package:" + packageName)

		outFileForServer, err := os.Create(outfileNameServer)
		if err != nil {
			panic(err)
		}

		outFileForClient, err := os.Create(outfileNameClient)
		if err != nil {
			panic(err)
		}

		srcFile, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(srcFile)
		if err != nil {
			panic(err)
		}
		srcFile.Close()
		srcFile, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}

		err = rpc.Generate(packageName, "", data, srcFile, outFileForServer, outFileForClient)
		if err != nil {
			panic(err)
		}
		srcFile.Close()
		outFileForServer.Close()
		outFileForClient.Close()

		err = exec.Command("go", "fmt", filename).Run()
		if err != nil {
			panic(err)
		}

		err = exec.Command("go", "fmt", outfileNameServer).Run()
		if err != nil {
			panic(err)
		}

		err = exec.Command("go", "fmt", outfileNameClient).Run()
		if err != nil {
			panic(err)
		}
	}
}
