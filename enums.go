package rewire

// AgentKeyStatus indicates an agents registration on a transport server
type AgentKeyStatus int

const (
	// AgentKeyError indicates an error reading the key status
	AgentKeyError = iota

	// AgentRegistered indicates an agent is registered
	AgentRegistered

	// AgentPending indicates the agent is waiting on approval
	AgentPending

	// AgentBlacklist indicates the agent key is blacklisted
	AgentBlacklist

	// AgentUnknown indicates the agent key is not known on this server
	AgentUnknown
)

func (s AgentKeyStatus) String() string {
	return []string{"KeyError", "KeyRegistered", "KeyPending", "KeyBlacklist", "KeyUnknown"}[s]
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

// ConnType defines agent or client
type ConnType int

const (
	// BadConnType is an unknown or failed connection type
	BadConnType = iota

	// AgentConnType is an agent connection
	AgentConnType

	// ClientConnType is a client connetion
	ClientConnType
)

func (s ConnType) String() string {
	return []string{"AgentConnection", "ClientConnection"}[s]
}

// ReqType is a request type identifier
type ReqType int

const (
	// InvalidReqType defines an unknown request type
	InvalidReqType = iota

	// EncryptedReqType is an encrypted payload
	EncryptedReqType

	// KeyReqType is a server key request
	KeyReqType

	// KeyRespType is a response to KeyReqType
	KeyRespType

	// AuthReqType is an authentication request
	AuthReqType

	// AuthRespType is a response to AuthReqType
	AuthRespType

	// OpenTunnReqType is an open tunnel request
	OpenTunnReqType

	// OpenTunnRespType is a response to OpenTunnReqType
	OpenTunnRespType
)

func (s ReqType) String() string {
	return []string{"InvalidReqType", "KeyReqType", "KeyRespType", "AuthReqType", "AuthRespType", "OpenTunnelReqType", "OpenTunneRespType"}[s]
}
