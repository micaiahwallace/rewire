package remoteagent

import (
	"crypto/rsa"
	"fmt"
	"io"
	"micaiahwallace/rewire"
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
	key, err := rwcrypto.Keys.GetPrivate("agent.key")
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

	// Start request to server
	newreq := &rewire.Request{Type: rewire.AgentConnType}
	struc.Pack(conn, newreq)

	// Create authentication signature
	sigstr, err := rwcrypto.CreateSignature(agent.key, "authenticate")
	if err != nil {
		panic("Unable to create valud signature")
	}

	// Get pem format of key
	pubpem, err := rwcrypto.GetKeyBytes(&agent.key.PublicKey, false)
	if err != nil {
		panic("Unable to get public key pem bytes")
	}

	// Create auth request structure
	authReq := &rewire.AuthRequest{
		PubKey: pubpem,
		Sig:    sigstr,
		// AuthCode: "1234",
	}
	struc.Pack(conn, authReq)

	// Receive auth response
	authresp := &rewire.AuthResp{}
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

	// Create forward session
	if err := agent.ForwardConnection(destAddr); err != nil {
		fmt.Println("Unable to create forward stream session", err)
	}
}

// ForwardConnection muxes a tcp connection and joins inbound streams to the destination host via new outbound connection per stream
func (agent *Agent) ForwardConnection(destAddr string) error {

	// Create smux server to accept connections
	muxer, err := smux.Server(*agent.conn, nil)
	if err != nil {
		fmt.Println("Unable to mux stream")
		return err
	}

	// Loop accept new streams
	for {

		// Accept a new stream
		inconn, err := muxer.AcceptStream()
		if err != nil {
			fmt.Println("Cannot accept stream mux", err)
			continue
		}

		// Create the forward in a new go routine
		go func(inconn *smux.Stream) {

			// Dial to destination
			outconn, err := net.Dial("tcp4", destAddr)
			if err != nil {
				fmt.Println("Cannot dial outbound connection")
				return
			}

			// connect inbound and outbound connections
			io.Copy(inconn, outconn)
			io.Copy(outconn, inconn)
		}(inconn)
	}

	return nil
}
