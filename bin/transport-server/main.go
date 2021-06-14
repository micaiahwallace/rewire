package main

import (
	"log"
	"micaiahwallace/rewire/transportsrv"
	"time"
)

func main() {

	// Always try to restart after failure
	for {
		if _, err := transportsrv.NewServer("", "3000"); err != nil {
			log.Printf("TS error: %v\n", err.Error())
			log.Println("Restarting server in 3 seconds.")
			time.Sleep(time.Second * 3)
		}
	}
}
