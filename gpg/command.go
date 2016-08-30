package main

import (
	"errors"
	"flag"
	"fmt"
	"net/rpc"
	"os"
)

var dispatcher *commandDispatcher

func init() {
	dispatcher = &commandDispatcher{
		commands: []command{
			exclusiveCommands(
				&listPublicKeys{},
				&listSecretKeys{},
				&importKeys{},
			),
		},
	}
}

type command interface {
	addFlags()
	canRun() bool
	run(*rpc.Client) error //XXX no need for an interface - for now
}

func exclusiveCommands(commands ...command) command {
	return &exclusiveCommand{commands}
}

type exclusiveCommand struct {
	commands []command
}

func (cmd *exclusiveCommand) addFlags() {
	for _, c := range cmd.commands {
		c.addFlags()
	}
}

func (cmd *exclusiveCommand) canRun() bool {
	for _, c := range cmd.commands {
		if c.canRun() {
			return true
		}
	}

	return false
}

func (cmd *exclusiveCommand) checkExclusive() error {
	canRun := false
	for _, c := range cmd.commands {
		if canRun && c.canRun() {
			return errors.New("conflicting commands")
		}

		canRun = c.canRun()
	}

	return nil
}

func (cmd *exclusiveCommand) run(client *rpc.Client) error {
	err := cmd.checkExclusive()
	if err != nil {
		return err
	}

	for _, c := range cmd.commands {
		if c.canRun() {
			return c.run(client)
		}
	}

	return nil
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
		if !cmd.canRun() {
			continue
		}

		err := cmd.run(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gpg: %s\n", err)
			os.Exit(-1)
		}
	}
}
