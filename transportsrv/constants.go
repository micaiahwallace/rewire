package transportsrv

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
