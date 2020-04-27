package transportsrv

const (
	// AgentConnType is an agent connection
	AgentConnType = iota
	// ClientConnType is a client connetion
	ClientConnType
)

// ConnType defines agent or client
type ConnType int

func (s ConnType) String() string {
	return []string{"AgentConnection", "ClientConnection"}[s]
}

// NewRequest is a new request identifying connection type
type NewRequest struct {
	Type ConnType `struc:"int8"`
}

// AgentAuthRequest is an agent auth request to transport server
type AgentAuthRequest struct {
	PubSize  int    `struc:"int16,sizeof=PubKey"`
	PubKey   []byte `struc:"[]int32"`
	SigSize  int    `struc:"int16,sizeof=Sig"`
	Sig      string
	CodeSize int `struc:"int8,sizeof=AuthCode"`
	AuthCode string
}

// AgentAuthResp is a response to AgentAuthRequest
type AgentAuthResp struct {
	Authenticated bool
}

// ClientAuthRequest is a client authentication request to transport server
type ClientAuthRequest struct {
}

// ClientAuthResp is an authentication response to the client
type ClientAuthResp struct {
	Authenticated bool
}

// OpenTunnelRequest defines an inbound request to open a tunnel
type OpenTunnelRequest struct {
	HostSize int `struc:"int8,sizeof=Host"`
	Host     string
	Port     int
	Agent    string
}
