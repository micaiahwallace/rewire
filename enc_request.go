package rewire

import (
	"bytes"
	"crypto/rsa"
	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
)

// EncryptedRequest holds an encrypted request to the server
type EncryptedRequest struct {
	Type    ReqType `struc:"int8"`
	Size    int     `struc:"int64,sizeof=Payload"`
	Payload []byte
}

// EncryptRequest encrypts a request with a public key
func EncryptRequest(req interface{}, key *rsa.PublicKey) (*EncryptedRequest, *EncryptedRequest, error) {

	// convert request to bytes
	var buff bytes.Buffer
	struc.Pack(&buff, req)
	reqbytes := buff.Bytes()

	// encrypt the bytes
	encbytes, err := rwcrypto.EncryptBytes(reqbytes, key)
	if err != nil {
		return nil, nil, err
	}

	// create and return encrypted request structure
	encreq := &EncryptedRequest{
		Type:    req.(Request).Type,
		Payload: encbytes,
	}

	// create receive payload
	encres := &EncryptedRequest{}

	return encreq, encres, nil
}
