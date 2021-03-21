package rpc

import (
	"context"
	"time"

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
// try to connect to rpc server in 5 seconds
func Init() error {
	conn, err := grpc.Dial(config.RPCTarget, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	client = CSTIRPC.NewCSTIRpcServerClient(conn)
	return nil
}
