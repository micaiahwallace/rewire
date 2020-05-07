package rewire

import (
	"bytes"
	"crypto/rsa"
	"micaiahwallace/rewire/rwcrypto"
	"net"

	"github.com/lunixbochs/struc"
)

// Request is an interface for generic request properties
type Request struct {
	// ReqType defines the request type
	Type ReqType
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

// DecryptResponse decrypts a response using a private key
func DecryptResponse(encres *EncryptedRequest, key *rsa.PrivateKey, res interface{}) error {

	// decrypt the payload
	resBytes, err := rwcrypto.DecryptBytes(encres.Payload, key)
	if err != nil {
		return err
	}

	// unpack response from decrypted bytes
	struc.Unpack(bytes.NewBuffer(resBytes), res)

	return nil
}

// SendRequest sends a request to the server and load response into resp
func SendRequest(conn *net.Conn, req interface{}, resp interface{}) error {

	// send request to server
	struc.Pack(*conn, req)

	// receive response
	struc.Unpack(*conn, resp)
	return nil
}

// SendRequestEncrypted sends requests encrypted with server key, receives response encrypted with client key
func SendRequestEncrypted(conn *net.Conn, req interface{}, resp interface{}, clientKey *rsa.PrivateKey, serverKey *rsa.PublicKey) error {

	// create encrypted payload wrapper
	encreq, encres, err := EncryptRequest(req, serverKey)
	if err != nil {
		return err
	}

	// send request and receive encrypted response
	SendRequest(conn, encreq, encres)

	// decrypt response into resp interface
	err = DecryptResponse(encreq, clientKey, resp)
	if err != nil {
		return err
	}

	return nil
}

// EncryptedRequest holds an encrypted request to the server
type EncryptedRequest struct {
	Type    ReqType `struc:"int8"`
	Size    int     `struc:"int64,sizeof=Payload"`
	Payload []byte
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
