package transportsrv

import (
	"crypto/rsa"
	"fmt"
	"micaiahwallace/rewire/rwutils"
	"net"

	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
	"github.com/xtaci/smux"
)

var agent1pubkey, _ = rwcrypto.ParsePublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0XIUgGx8I/rcwAoekqc6
/wNyG/bEKeT3lQiDv5WP+G8tfh2rcggQj9CVof37qJ2lytuh1rqpLrvWV2dQFacK
Pnemvx0BVrKBY35O7QQCoV+j1BD/Hcf9AoijryUT9LnvZU92KPh1imT6QNGln2Uv
Gr4RSzmoLDzKTcZIrcp/Vfx1uLPXj2dkDgZ46O/OIb5MLQdqMagxJJe9CeC1GW0o
i9+T5ZnAgI073DEfNjfXUxGdhZ4AOFC5IXcf2fu4aDabS9jm3LLWEwOba+0Lf6Xh
3CsQpw6zRjuYRB2ys7mZ2bmxOWWxK/zpCj4EaMDiAVar0oiB9Q43CoggTywUAdnE
xQIDAQAB
-----END PUBLIC KEY-----`))

var (
	registeredAgents []Agent = []Agent{
		{
			Pubkey: agent1pubkey,
		},
	}
)

// Server acts as conduit between remote agent and client
type Server struct {

	// underlying server tcp socket
	socket net.Listener

	// rsa keypair
	key *rsa.PrivateKey

	// active agent connections
	agents map[string]*AgentConn

	// active client connections
	clients []string

	// quit channel
	quit chan int
}

// New creates a new transport server
func New() *Server {

	// Create new server
	server := Server{
		quit:   make(chan int),
		agents: make(map[string]*AgentConn),
	}

	// Create necessary directory structure
	errs := rwutils.MakePaths([]string{
		"keys/agent/pending",
		"keys/agent/accepted",
		"keys/agent/blacklist",
		"keys/client/pending",
		"keys/client/accepted",
		"keys/client/blacklist",
	})
	if len(errs) > 0 {
		fmt.Println(errs)
		panic("Unable to create directory structures.")
	}

	return &server
}

// Start the transport server listener
func (server *Server) Start(port int) error {
	quitCheck := false

	// create and start listener
	addr := fmt.Sprintf(":%d", port)
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Unable to bind socket")
		return err
	}
	fmt.Println("Transport server running", addr)

	// load server keypair
	key, err := rwcrypto.LoadKeypair("server.key")
	server.key = key
	if err != nil {
		fmt.Println("Unable to load or generate server keys")
		return err
	}

	// loop and accept connections
	for {
		select {

		// stop the server channel check
		case <-server.quit:
			quitCheck = true

		// accept a connection
		default:
			conn, err := socket.Accept()
			if err != nil {
				fmt.Println("Unable to accept connection")
			} else {
				go server.handleConnection(conn)
			}
		}

		// quit the server
		if quitCheck {
			fmt.Println("Closing server")
			break
		}
	}

	return nil
}

// Stop the server listener
func (server *Server) Stop() {
	server.quit <- 0
}

// Handle agent authentication flow
func (server *Server) handleAgent(conn net.Conn) {

	// Get agent auth request
	agentAuthReq := &AgentAuthRequest{}
	struc.Unpack(conn, agentAuthReq)

	// decode the public key
	pubkey, err := rwcrypto.ParsePublicKey(agentAuthReq.PubKey)
	if err != nil {
		fmt.Println("Unable to parse public key", err)

		// send agent auth response
		agentAuthResponse := &AgentAuthResp{
			Authenticated: false,
		}
		struc.Pack(conn, agentAuthResponse)
		return
	}

	// verify public key matches signature
	validSig := rwcrypto.ValidateSignature(pubkey, agentAuthReq.Sig, "authenticate")
	if !validSig {
		fmt.Println("Signature does not match public key")

		// send agent auth response
		agentAuthResponse := &AgentAuthResp{
			Authenticated: false,
		}
		struc.Pack(conn, agentAuthResponse)
		return
	}

	// create new agent connection
	newAgent := NewAgentConn(pubkey, &conn)

	// check if public key is already registered
	regStatus := newAgent.Agent.KeyStatus()

	if regStatus == AgentRegistered {
		// store agent connection by public key
		server.agents[newAgent.Agent.ID] = newAgent
		fmt.Println("registered agent connected", newAgent.Agent.ID)

		// send agent auth response
		agentAuthResponse := &AgentAuthResp{
			Authenticated: true,
		}
		struc.Pack(conn, agentAuthResponse)
		return
	}

	if regStatus == AgentUnknown {
		fmt.Println("agent not yet registered", newAgent.Agent.ID)

		// Register agent as pending
		if err = newAgent.Agent.SetKeyStatus(AgentPending); err != nil {
			fmt.Println("unable to register pending agent", err, newAgent.Agent.ID)
		}

	}

	// send agent auth response
	agentAuthResponse := &AgentAuthResp{
		Authenticated: false,
	}
	struc.Pack(conn, agentAuthResponse)
}

// Handle client authentication flow
func (server *Server) handleClient(conn net.Conn) *smux.Session {

	// multiplex tcp stream
	session, err := smux.Server(conn, nil)
	if err != nil {
		panic(err)
	}

	return session
}

// Logic for processing connections
func (server *Server) handleConnection(conn net.Conn) {

	// identify connection type, agent or client
	inReq := &NewRequest{}
	struc.Unpack(conn, inReq)
	fmt.Println("Type:", inReq.Type)

	// Handle connection logic based on connection type
	switch inReq.Type {
	case AgentConnType:
		server.handleAgent(conn)
	case ClientConnType:
		server.handleClient(conn)
	default:
		fmt.Println("Unknown connection type")
	}
}
