package rewire

// ConfigTemplate holds the configuration format for the rewire lib
type ConfigTemplate struct {

	// root path to pki directory
	PKIRoot string

	// path to client key file
	LocalKey string

	// path to server key file
	ServerKey string

	// Key length in bits when generating
	KeyBitLength int

	// string to sign when authenticating clients
	AuthSigStr string
}
