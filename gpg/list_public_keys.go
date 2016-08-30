package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/rpc"
	"os"

	"golang.org/x/crypto/openpgp"

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

func (cmd *listPublicKeys) canRun() bool {
	return cmd.active
}

func (cmd *listPublicKeys) run(c *rpc.Client) error {
	if !cmd.active {
		return nil
	}

	var pubKeyRing types.GetKeyRingResult
	err := c.Call("Protocol.GetPublicKeyRing", types.VoidArg{}, &pubKeyRing)
	if err != nil {
		return err
	}

	for _, k := range pubKeyRing.Keys {
		b := bytes.NewBufferString(k.KeyData)
		entities, err := openpgp.ReadArmoredKeyRing(b)
		if err != nil {
			// Errors are like:
			//gpg: conversion from `utf-8' to `US-ASCII' failed: Illegal byte sequence
			fmt.Println("gpg:", err)
			continue
		}

		publicKeyringFormat(os.Stdout, entities)
	}

	return nil
}
