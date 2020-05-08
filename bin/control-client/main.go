package main

import (
	"fmt"
	"micaiahwallace/rewire/ctrlclient"
)

func main() {
	if err := ctrlclient.Run("localhost", "3000"); err != nil {
		fmt.Println("Control Client Error:", err.Error())
	}
}
