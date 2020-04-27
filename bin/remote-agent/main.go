package main

import "micaiahwallace/rewire/remoteagent"

func main() {
	agent := remoteagent.New()
	agent.Connect("localhost", 3000)
}
