package main

import (
	"micaiahwallace/rewire/transportsrv"
)

func main() {
	server := transportsrv.New()
	server.Start(3000)
}
