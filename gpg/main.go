package main

import (
	"log"
	"os"
)

const (
	defaultDaemonAddress = "/var/run/nyms/nyms.socket"
)

func main() {
	//XXX Should we connect before running any command?
	client, err := connect(defaultDaemonAddress)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	dispatcher.dispatch(client)
}
