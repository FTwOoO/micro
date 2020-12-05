package main

import (
	"time"
	"tracingtest/examples"
)

func main() {
	examples.HttpTest()
	examples.GrpcTest()
	time.Sleep(30 * time.Second)
}
