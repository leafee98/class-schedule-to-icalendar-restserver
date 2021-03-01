package rpc

import (
	"context"

	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/rpc/CSTIRPC"
	"google.golang.org/grpc"
)

var client CSTIRPC.CSTIRpcServerClient

// JSONGenerate generate the icalendar result, request string and return string and error
func JSONGenerate(content string) (string, error) {
	res, err := client.JsonGenerate(context.TODO(), &CSTIRPC.ConfJson{Content: content})
	return res.GetContent(), err
}

// Init initialize RPCClient object to use rpc of rpcserver
func Init() error {
	conn, err := grpc.Dial(config.RPCTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	client = CSTIRPC.NewCSTIRpcServerClient(conn)
	return nil
}
