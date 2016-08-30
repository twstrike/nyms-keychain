package main

import (
	"flag"
	"net/rpc"
)

var dispatcher *commandDispatcher

func init() {
	dispatcher = &commandDispatcher{
		commands: []command{
			&listPublicKeys{},
		},
	}
}

type command interface {
	addFlags()
	run(*rpc.Client) //XXX no need for an interface - for now
}

type commandDispatcher struct {
	commands []command
}

func (c *commandDispatcher) dispatch(client *rpc.Client) {
	for _, cmd := range c.commands {
		cmd.addFlags()
	}

	flag.Parse()

	for _, cmd := range c.commands {
		cmd.run(client)
	}
}
