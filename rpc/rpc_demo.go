package rpc

import "context"

type HelloRequest struct {
	Send int `bson:"send" json:"send"`
}
type HelloResponse struct {
	Ok bool `bson:"ok" json:"ok"`
}
type DemoSevice interface {
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
}
