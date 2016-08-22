package gui

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/twstrike/nyms-agent/protocol/types"
)

type client struct {
	rpc *rpc.Client
}

func (c *client) establishConn() {
	conn, err := net.Dial("unix", "/tmp/nyms.sock")
	if err != nil {
		log.Fatal(err)
	}

	c.rpc = jsonrpc.NewClient(conn)
}

func (c *client) pubKeyRing() types.GetKeyRingResult {
	c.establishConn()
	var pubKeyRing types.GetKeyRingResult
	err := c.rpc.Call("Protocol.GetPublicKeyRing", types.VoidArg{}, &pubKeyRing)
	if err != nil {
		log.Fatal("GetPublicKeyRing error:", err)
	}
	return pubKeyRing
}

func (c *client) createKeyPair(name, email, comment string) types.GetKeyInfoResult {
	c.establishConn()
	var keyinfo types.GetKeyInfoResult
	err := c.rpc.Call("Protocol.GenerateKeys", types.GenerateKeysArgs{name, email, comment, ""}, &keyinfo)
	if err != nil {
		log.Fatal("GenerateKeys error:", err)
	}
	return keyinfo
}
