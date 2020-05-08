package rewire

import (
	"crypto/rsa"
	"micaiahwallace/rewire/rwcrypto"
)

// KeyRequest is a request structure for server key request
type KeyRequest struct {
	Type ReqType
}

// KeyResp is a response structure for a key request
type KeyResp struct {
	KeySize int    `struc:"int32,sizeof=Key"`
	Key     []byte `struc:"[]int32"`
}

// ServerKeyRequest creates a server key request
func ServerKeyRequest() (*KeyRequest, *KeyResp) {
	req := &KeyRequest{Type: KeyReqType}
	res := &KeyResp{}
	return req, res
}

// ServerKeyResponse creates a server key response
func ServerKeyResponse(key *rsa.PublicKey) (*KeyResp, error) {

	// convert key to bytes
	keybytes, err := rwcrypto.GetKeyBytes(key, false)
	if err != nil {
		return nil, err
	}

	// create response struct
	res := &KeyResp{Key: keybytes}
	return res, nil
}
