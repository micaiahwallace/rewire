package main

import "micaiahwallace/rewire/remoteagent"

func main() {
	if err := remoteagent.Run("localhost", "3000"); err != nil {
		panic(err)
	}
}
