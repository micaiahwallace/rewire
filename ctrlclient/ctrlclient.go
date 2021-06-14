package ctrlclient

import (
	"micaiahwallace/rewire"
)

// Run runs a control client session
func Run(host, port string) error {

	// Create rewire client with control client type
	_, err := rewire.InitClient(rewire.ClientConnType, host, port)
	if err != nil {
		return err
	}

	return nil
}
