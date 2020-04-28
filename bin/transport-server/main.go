package main

import (
	"micaiahwallace/rewire/transportsrv"
)

func main() {
	server := transportsrv.New()
	if err := server.Start("3000"); err != nil {
		panic(err)
	}
}
