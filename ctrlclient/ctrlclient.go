package ctrlclient

import (
	"micaiahwallace/rewire"
)

// Run runs a control client session
func Run(host, port string) error {

	// Create rewire client
	client, err := rewire.InitClient(rewire.ClientConnType, host, port)
	if err != nil {
		return err
	}

	// Authenticate to server
	if err := client.Authenticate(); err != nil {
		return err
	}

	return nil
}
