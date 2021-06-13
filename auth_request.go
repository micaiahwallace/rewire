package rewire

import (
	"crypto/rsa"
	"errors"
	"micaiahwallace/rewire/rwcrypto"
)

// AuthRequest is a request to authenticate a connection with the server
type AuthRequest struct {
	Type     ReqType  `struc:"int8"`
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
func CreateAuthRequest(connType ConnType, key *rsa.PrivateKey) (*AuthRequest, *AuthResp, error) {

	// Create auth signature
	sig, err := rwcrypto.CreateSignature(key, Config.AuthSigStr)
	if err != nil {
		return nil, nil, errors.New("cannot sign authentication message")
	}

	// Get key bytes
	keybytes, err := rwcrypto.GetKeyBytes(key.PublicKey, false)
	if err != nil {
		return nil, nil, errors.New("cannot convert key to bytes")
	}

	// create auth request
	authreq := &AuthRequest{
		Type:     AuthReqType,
		ConnType: connType,
		PubKey:   keybytes,
		Sig:      sig,
	}

	return authreq, &AuthResp{}, nil
}
