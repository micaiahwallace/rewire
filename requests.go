package rewire

import (
	"bytes"
	"crypto/rsa"
	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
)

// Request is an interface for generic request properties
type Request struct {
	// ReqType defines the request type
	Type ReqType
}

// Request sends a request and recieves a result from the server
func (client *Client) Request(req, resp interface{}) error {

	// pack request over tcp connection
	if err := struc.Pack(*client.conn, req); err != nil {
		return err
	}

	// unpack response from tcp connection
	if err := struc.Unpack(*client.conn, resp); err != nil {
		return err
	}

	return nil
}

// DecryptResponse decrypts a response using a private key
func (client *Client) DecryptResponse(encres *EncryptedRequest, key *rsa.PrivateKey, res interface{}) error {

	// decrypt the payload
	resBytes, err := rwcrypto.DecryptBytes(encres.Payload, key)
	if err != nil {
		return err
	}

	// unpack response from decrypted bytes
	if err := struc.Unpack(bytes.NewBuffer(resBytes), res); err != nil {
		return err
	}

	return nil
}

// SendRequestEncrypted sends requests encrypted with server key, receives response encrypted with client key
func (client *Client) SendRequestEncrypted(req interface{}, resp interface{}) error {

	// Get local key
	localKey, err := rwcrypto.Keys.GetPrivate(Config.LocalKey)
	if err != nil {
		return err
	}

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
	if err := client.Request(encreq, encres); err != nil {
		return err
	}

	// decrypt response into resp interface
	if err := client.DecryptResponse(encres, localKey, resp); err != nil {
		return err
	}

	return nil
}

// OpenTunnelRequest defines an inbound request to open a tunnel
type OpenTunnelRequest struct {
	AgentSize int `struc:"int16,sizeof=Agent"`
	Agent     string
	HostSize  int `struc:"int8,sizeof=Host"`
	Host      string
	PortSize  int `struc:"int8,sizeof=Port"`
	Port      string
}
