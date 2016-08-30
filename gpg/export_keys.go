package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"

	"github.com/twstrike/nyms-agent/protocol/types"
)

type exportPublicKeys struct {
	active        bool
	armoredOutput bool
}

func (cmd *exportPublicKeys) addFlags() {
	flag.BoolVar(&cmd.active, "export", false, "export keys")
	flag.BoolVar(&cmd.armoredOutput, "armor", false, "create ascii armored output")
}

func (cmd *exportPublicKeys) canRun() bool {
	return cmd.active
}

func (cmd *exportPublicKeys) run(c *rpc.Client) error {
	if !cmd.active {
		return nil
	}

	ret := new(types.ExportEntitiesResult)
	err := c.Call("Protocol.ExportEntities", types.ExportEntities{
		ArmoredOutput: cmd.armoredOutput,
	}, &ret)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(ret.Output)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stdout, "")
	return err
}
