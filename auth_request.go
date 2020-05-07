package rewire

import (
	"crypto/rsa"
	"errors"
	"micaiahwallace/rewire/rwcrypto"
)

// AuthRequest is a request to authenticate a connection with the server
type AuthRequest struct {
	ConnType ConnType `struc:"int8"`
	PubSize  int      `struc:"int16,sizeof=PubKey"`
	PubKey   []byte   `struc:"[]int32"`
	SigSize  int      `struc:"int16,sizeof=Sig"`
	Sig      string
}

// AuthResp is a response to AuthRequest
type AuthResp struct {
	Authenticated bool `struc:"bool"`
	MessageSize   int  `struc:"int32,sizeof=Message"`
	Message       string
}

// CreateAuthRequest creates an auth request for a connection type and public key
func CreateAuthRequest(t ConnType, key *rsa.PrivateKey) (*AuthRequest, error) {

	// Create auth signature
	sig, err := rwcrypto.CreateSignature(key, "authenticate")
	if err != nil {
		return nil, errors.New("Cannot sign authentication message")
	}

	// Get key bytes
	keybytes, err := rwcrypto.ExportKeyBytes(key.PublicKey, false)
	if err != nil {
		return nil, errors.New("Cannot convert key to bytes")
	}

	// create auth request
	authreq := &AuthRequest{
		ConnType: t,
		PubKey:   keybytes,
		Sig:      sig,
	}

	return authreq, nil
}
