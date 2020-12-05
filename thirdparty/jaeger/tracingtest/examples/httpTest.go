package examples

import (
	"flag"
	"fmt"
)

var (
	serverPort = flag.String("port", "8000", "server port")
)

const (
	server = "httpServer"
	client = "httpClient"
)

func HttpTest() {

	go startServer()

	runClient(NewJaegerTracer(client))
	fmt.Println("http done")

}

func startServer() {
	runServer(NewJaegerTracer(server))
}
