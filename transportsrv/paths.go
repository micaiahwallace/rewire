package transportsrv

import (
	"path"
)

const (

	// KeyDir root key directory
	KeyDir = "keys"

	// AcceptedKeyDir name for accepted key directory
	AcceptedKeyDir = "accepted"

	// PendingKeyDir name for accepted key directory
	PendingKeyDir = "pending"

	// BlacklistKeyDir name for accepted key directory
	BlacklistKeyDir = "blacklist"

	// ServerKeyName name of server key file
	ServerKeyName = "server.key"
)

var (

	// AgentKeyDir root agent keys
	AgentKeyDir = path.Join(KeyDir, "agent")

	// AcceptedAgentKeyDir accepted agent keys
	AcceptedAgentKeyDir = path.Join(AgentKeyDir, AcceptedKeyDir)

	// PendingAgentKeyDir pending agent keys
	PendingAgentKeyDir = path.Join(AgentKeyDir, PendingKeyDir)

	// BlacklistAgentKeyDir blacklited agent keys
	BlacklistAgentKeyDir = path.Join(AgentKeyDir, BlacklistKeyDir)

	// ClientKeyDir root client keys
	ClientKeyDir = path.Join(KeyDir, "client")

	// AcceptedClientKeyDir accepted agent keys
	AcceptedClientKeyDir = path.Join(ClientKeyDir, AcceptedKeyDir)

	// PendingClientKeyDir pending agent keys
	PendingClientKeyDir = path.Join(ClientKeyDir, PendingKeyDir)

	// BlacklistClientKeyDir blacklited agent keys
	BlacklistClientKeyDir = path.Join(ClientKeyDir, BlacklistKeyDir)
)

// KeyPath creates a file path string to a given key
func KeyPath(keydir, keyname string) string {
	return path.Join(keydir, keyname)
}
