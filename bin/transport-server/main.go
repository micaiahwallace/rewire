package main

import (
	"micaiahwallace/rewire/transportsrv"
)

func main() {
	if _, err := transportsrv.NewServer("", "3000"); err != nil {
		panic(err)
	}
}
