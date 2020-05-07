package main

import "micaiahwallace/rewire/ctrlclient"

func main() {
	client := ctrlclient.New()
	client.Connect("localhost:3000")
}
