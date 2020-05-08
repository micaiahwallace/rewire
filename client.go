package rewire

import (
	"errors"
	"micaiahwallace/rewire/rwcrypto"
	"net"
)

// Client holds client state and connection
type Client struct {
	Type ConnType
	Host string
	Port string
	conn *net.Conn
}

// Config store configuration for rewire
var Config *ConfigTemplate

// InitClient creates and sets up a new client
func InitClient(clientType ConnType, host, port string) (*Client, error) {

	// validate host param
	if host == "" {
		return nil, errors.New("Host must be defined")
	}

	// validate port param
	if port == "" {
		return nil, errors.New("Port must be defined")
	}

	// Init rewire lib
	InitLib()

	// Create new client
	client := Client{Host: host, Port: port, Type: clientType}

	// Initiate server connection
	err := client.Connect()
	if err != nil {
		return nil, err
	}

	// Setup keystore keys
	err = client.SetupKeys()
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Connect to a transport server
func (client *Client) Connect() error {

	// Dial connection to server
	conn, err := net.Dial("tcp4", net.JoinHostPort(client.Host, client.Port))
	if err != nil {
		return err
	}

	// store socket reference
	client.conn = &conn

	return nil
}

// SetupKeys ensures the client and server keys exist
func (client *Client) SetupKeys() error {

	// Check if local key exists
	if !rwcrypto.Keys.Contains(Config.LocalKey) {

		// generate new local key
		newKey, err := rwcrypto.GenerateKey(2048)
		if err != nil {
			return err
		}

		// save to the keystore
		rwcrypto.Keys.SavePrivate(Config.LocalKey, newKey)
	}

	// Check if server key in key store
	if !rwcrypto.Keys.Contains(Config.ServerKey) {

		// create and send key request
		serverKeyReq, serverKeyResp := ServerKeyRequest()
		err := client.Request(serverKeyReq, &serverKeyResp)
		if err != nil {
			return err
		}

		// Convert key bytes to rsa key
		serverKey, err := rwcrypto.ParsePublicPEM(serverKeyResp.Key)
		if err != nil {
			return err
		}

		// save resulting key to keystore
		rwcrypto.Keys.SavePublic(Config.ServerKey, serverKey)
	}

	return nil
}

// Authenticate against the server via client key
func (client *Client) Authenticate() error {

	// Get local key
	localKey, err := rwcrypto.Keys.GetPrivate(Config.LocalKey)
	if err != nil {
		return err
	}

	// Create auth request payload
	authreq, authresp, err := CreateAuthRequest(client.Type, localKey)
	if err != nil {
		return err
	}

	// Send auth request to server
	if err := client.SendRequestEncrypted(&authreq, &authresp); err != nil {
		return err
	}

	// return error if not authenticated
	if !authresp.Authenticated {
		return errors.New("Client not authenticated")
	}

	return nil
}
