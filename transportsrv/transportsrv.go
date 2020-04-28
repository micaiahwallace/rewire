package transportsrv

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"micaiahwallace/rewire/rwutils"
	"net"

	"micaiahwallace/rewire/rwcrypto"

	"github.com/lunixbochs/struc"
	"github.com/xtaci/smux"
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

	// Initialize server dependencies
	err := server.Setup()
	if err != nil {
		panic("Unable to setup server dependencies.")
	}

	return &server
}

// Setup required directories and objects
func (server *Server) Setup() error {

	// Create pki directory structure
	errs := rwutils.MakePaths([]string{
		AcceptedAgentKeyDir,
		PendingAgentKeyDir,
		BlacklistAgentKeyDir,
		AcceptedClientKeyDir,
		PendingClientKeyDir,
		BlacklistClientKeyDir,
	})

	// Check for errors creating directories
	if len(errs) > 0 {
		return errors.New("Unable to provision pki dir tree")
	}

	// load server keypair
	key, err := rwcrypto.LoadKeypair(ServerKeyName)
	if err != nil {
		return errors.New("Unable to load or generate server keys")
	}
	server.key = key

	return nil
}

// Start the transport server listener
func (server *Server) Start(port string) error {
	quitCheck := false

	// create and start listener
	addr := net.JoinHostPort("", port)
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server.socket = socket
	fmt.Println("Transport server running", addr)

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

	switch regStatus {

	// Registered agent will stay connected
	case AgentRegistered:
		// store agent connection by public key
		server.agents[newAgent.Agent.ID] = newAgent
		fmt.Println("registered agent connected", newAgent.Agent.ID)

		// send agent auth response
		agentAuthResponse := &AgentAuthResp{
			Authenticated: true,
		}
		struc.Pack(conn, agentAuthResponse)
		return

	// Unknown agent will become pending and disconnect
	case AgentUnknown:
	case AgentPending:
		fmt.Println("agent not yet registered", newAgent.Agent.ID)

		// Register agent as pending
		if err = newAgent.Agent.SetKeyStatus(AgentPending); err != nil {
			fmt.Println("unable to register pending agent", err, newAgent.Agent.ID)
		}

	// Blacklit agent will disconnect
	case AgentBlacklist:
		fmt.Println("agent key blacklisted", newAgent.Agent.ID)

	default:
		fmt.Println("agent key error", newAgent.Agent.ID)
	}

	// Send not authenticated
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
	fmt.Println("agents", server.agents)
}
