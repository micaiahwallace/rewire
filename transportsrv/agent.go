package transportsrv

import (
	"crypto/rsa"
	"net"

	"github.com/xtaci/smux"
)

// Agent is a remote agent connected to the transport server
type Agent struct {

	// registered agent public key
	Pubkey *rsa.PublicKey
}

// AgentConnStatus is the current status of the connection
type AgentConnStatus int

const (

	// AgentDown lost connection
	AgentDown = iota

	// AgentUp connection active but not authenticated
	AgentUp

	// AgentAuth connected and authenticated
	AgentAuth
)

func (s AgentConnStatus) String() string {
	return []string{"AgentDown", "AgentUp", "AgentAuth"}[s]
}

// AgentConn represents an active agent connection
type AgentConn struct {

	// agent details
	Agent *Agent

	// agent net connection
	Conn *net.Conn

	// connection status
	Status AgentConnStatus

	// connected agent tcp smux stream
	Session *smux.Session
}

// NewAgentConn returns an active agent connection
func NewAgentConn(key *rsa.PublicKey, socket *net.Conn) *AgentConn {
	agent := Agent{Pubkey: key}
	ac := AgentConn{
		Agent:  &agent,
		Status: AgentAuth,
		Conn:   socket,
	}
	return &ac
}
