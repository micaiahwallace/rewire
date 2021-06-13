package rewire

import (
	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
)

// SendServerRequest sends a request and recieves a result from the server
func (client *Client) SendServerRequest(req, resp interface{}) error {

	// pack request over net connection
	if err := struc.Pack(*client.conn, req); err != nil {
		return err
	}

	// unpack response from net connection
	if err := struc.Unpack(*client.conn, resp); err != nil {
		return err
	}

	return nil
}

// SendRequestEncrypted sends requests encrypted with server key, receives response encrypted with client key
func (client *Client) SendRequestEncrypted(req interface{}, resp interface{}) error {

	// Get server key
	serverKey, err := rwcrypto.Keys.GetPublic(Config.ServerKey)
	if err != nil {
		return err
	}

	// create encrypted payload wrapper
	encreq, encres, err := EncryptRequest(req, serverKey)
	if err != nil {
		return err
	}

	// send request and receive encrypted response
	if err := client.SendServerRequest(encreq, encres); err != nil {
		return err
	}

	// decrypt response into resp interface
	if err := DecryptRequest(encres, resp); err != nil {
		return err
	}

	return nil
}
