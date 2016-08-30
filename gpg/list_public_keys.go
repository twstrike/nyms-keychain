package main

import (
	"flag"
	"log"
	"net/rpc"
	"os"

	"github.com/twstrike/nyms-agent/protocol/types"
)

type listPublicKeys struct {
	active bool
}

func (cmd *listPublicKeys) addFlags() {
	//XXX Should we add an alias for
	//--list-keys, -k
	flag.BoolVar(&cmd.active, "list-public-keys", false, "list keys")
}

func (cmd *listPublicKeys) run(c *rpc.Client) {
	if !cmd.active {
		return
	}

	var pubKeyRing types.GetKeyRingResult
	err := c.Call("Protocol.GetPublicKeyRing", types.VoidArg{}, &pubKeyRing)
	if err != nil {
		log.Fatal("GetPublicKeyRing error:", err)
	}

	publicKeyringFormat(os.Stdout, pubKeyRing.Keys)
}
