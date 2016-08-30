package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/rpc"
	"os"

	"github.com/twstrike/nyms-agent/protocol/types"
	"golang.org/x/crypto/openpgp"
)

type listSecretKeys struct {
	active bool
}

func (cmd *listSecretKeys) addFlags() {
	//XXX Should we add an alias for -K ?
	flag.BoolVar(&cmd.active, "list-secret-keys", false, "list secret keys")
}

func (cmd *listSecretKeys) canRun() bool {
	return cmd.active
}

func (cmd *listSecretKeys) run(c *rpc.Client) error {
	if !cmd.active {
		return nil
	}

	var secKeyRing types.GetKeyRingResult
	err := c.Call("Protocol.GetSecretKeyRing", types.VoidArg{}, &secKeyRing)
	if err != nil {
		return err
	}

	for _, k := range secKeyRing.Keys {
		b := bytes.NewBufferString(k.KeyData)
		entities, err := openpgp.ReadArmoredKeyRing(b)
		if err != nil {
			// Errors are like:
			//gpg: conversion from `utf-8' to `US-ASCII' failed: Illegal byte sequence
			fmt.Println("gpg:", err)
			continue
		}

		//XXX it is the same as --list-public-keys,
		//but the output also have "ssb" and "ssb>" lines.
		publicKeyringFormat(os.Stdout, entities)
	}

	return nil
}
