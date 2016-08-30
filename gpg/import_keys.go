package main

import (
	"bytes"
	"flag"
	"io"
	"net/rpc"
	"os"

	"github.com/twstrike/nyms-agent/protocol/types"
)

type importKeys struct {
	active bool
}

func (cmd *importKeys) addFlags() {
	flag.BoolVar(&cmd.active, "import", false, "import/merge keys")
}

func (cmd *importKeys) canRun() bool {
	return cmd.active
}

func (cmd *importKeys) run(c *rpc.Client) error {
	if !cmd.active {
		return nil
	}

	b := new(bytes.Buffer)
	_, err := io.Copy(b, os.Stdin)
	if err != nil {
		return err
	}

	var void types.VoidArg
	return c.Call("Protocol.ImportEntities", types.ImportEntities{
		ArmoredEntities: b.String(),
	}, &void)
}
