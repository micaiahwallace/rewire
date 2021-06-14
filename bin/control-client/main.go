package main

import (
	"fmt"
	"micaiahwallace/rewire/ctrlclient"
	"time"
)

func main() {

	// Always try to reconnect
	for {

		// Attempt connection
		if err := ctrlclient.Run("localhost", "3000"); err != nil {
			fmt.Println("Control Client Error:", err.Error())
			fmt.Println("Retrying connection in 3 seconds.")
			time.Sleep(time.Second * 3)
		}
	}
}
