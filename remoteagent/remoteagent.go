package remoteagent

import (
	"crypto/rsa"
	"fmt"
	"micaiahwallace/rewire/rwcrypto"
	"micaiahwallace/rewire/transportsrv"
	"net"

	"github.com/xtaci/smux"

	"github.com/lunixbochs/struc"
)

// Agent is a remote agent to the transport server
type Agent struct {

	// private key for agent
	key *rsa.PrivateKey

	// connection to server
	conn *net.Conn
}

// New creates a new agent
func New() *Agent {
	key, err := rwcrypto.LoadKeypair("agent.key")
	if err != nil {
		panic("Unable to create keypair")
	}
	agent := Agent{
		key: key,
	}
	return &agent
}

// Connect establishes a tcp socket with a transport server
func (agent *Agent) Connect(host string, port int) {

	// Dial connection to server
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic("Can't connect to server")
	}
	agent.conn = &conn

	// Start request
	newreq := &transportsrv.NewRequest{Type: transportsrv.AgentConnType}
	struc.Pack(conn, newreq)

	// Start authentication
	sigstr, err := rwcrypto.CreateSignature(agent.key, "authenticate")
	if err != nil {
		panic("Unable to create valud signature")
	}
	pubpem, err := rwcrypto.ExportKeyBytes(&agent.key.PublicKey, false)
	if err != nil {
		panic("Unable to get public key pem bytes")
	}
	authReq := &transportsrv.AgentAuthRequest{
		PubKey:   pubpem,
		Sig:      sigstr,
		AuthCode: "1234",
	}
	struc.Pack(conn, authReq)

	// Receive auth response
	authresp := &transportsrv.AgentAuthResp{}
	struc.Unpack(conn, authresp)
	fmt.Println("authenticated:", authresp.Authenticated)

	if !authresp.Authenticated {
		return
	}

	// Receive connection request
	openReq := &transportsrv.OpenTunnelRequest{}
	struc.Unpack(conn, openReq)
	destAddr := net.JoinHostPort(openReq.Host, openReq.Port)
	fmt.Println("open tunnel to", openReq.Host, openReq.Port)

	// Filter out unauthorized requests

	// Multiplex connection and send to the open request destination
	muxer, err := smux.Server(conn, smux.DefaultConfig())
	if err != nil {
		panic("Unable to mux stream")
	}
	for {
		inconn, err := muxer.AcceptStream()
		if err != nil {
			panic("Stream receive error")
		}
		// dial to destAddr
	}
}
