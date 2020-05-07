package ctrlclient

import (
	"crypto/rsa"
	"micaiahwallace/rewire"
	"micaiahwallace/rewire/rwcrypto"
	"net"
)

// Client a transport server control client
type Client struct {
	conn      *net.Conn
	key       *rsa.PrivateKey
	serverKey *rsa.PublicKey
}

// New creates a new control client
func New() *Client {
	client := Client{}

	// Load keypair into client
	key, err := rwcrypto.LoadKeypair("client.key")
	if err != nil {
		panic("Cannot load private key")

	}
	client.key = key

	return &client
}

// Connect to a transport server
func (client *Client) Connect(addr string) error {

	// Create tcp connection to server
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		panic(err)
	}
	client.conn = &conn

	// Retrieve server public key
	err = client.FetchPublicKey()
	if err != nil {
		panic(err)
	}

	// Authenticate to the server
	err = client.Authenticate()
	if err != nil {
		panic(err)
	}

	return nil
}

// FetchPublicKey stores the server key if not already present
func (client *Client) FetchPublicKey() error {

	// Create and send key request
	keyreq, keyresp := rewire.CreateKeyRequest(client.key.PublicKey)
	rewire.SendRequest(client.conn, keyreq, keyresp)
}

// Authenticate the control client
func (client *Client) Authenticate() error {

	// Create auth request
	authreq, err := rewire.CreateAuthRequest(rewire.ClientConnType, client.key)
	if err != nil {
		return err
	}

	// send auth request to server and get response
	authresp := &rewire.AuthResp{}
	// rewire.EncryptRequest(authresp, key.PublicKey)
	rewire.SendRequest(client.conn, authreq, authresp)

	// send user defined command to server

	return nil
}
