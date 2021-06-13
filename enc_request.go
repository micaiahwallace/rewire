package rewire

import (
	"bytes"
	"crypto/rsa"
	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
)

// EncryptRequest encrypts a request with a public key
func EncryptRequest(req interface{}, key *rsa.PublicKey) (*Request, *Request, error) {

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
	encreq := &Request{
		Type:      req.(Request).Type,
		Encrypted: true,
		Payload:   encbytes,
	}

	// create receive payload
	encres := &Request{}

	return encreq, encres, nil
}

// DecryptRequest decrypts a request using a private key
func DecryptRequest(encres *Request, res interface{}) error {

	// Get server key
	key, keyErr := rwcrypto.Keys.GetPrivate(Config.LocalKey)
	if keyErr != nil {
		return keyErr
	}

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
