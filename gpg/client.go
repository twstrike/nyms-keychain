package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func connect(addr string) (*rpc.Client, error) {
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(conn), nil
}
