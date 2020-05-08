package rewire

import "micaiahwallace/rewire/rwcrypto"

// InitLib inits shared rewire logic
func InitLib() {

	// load config from cli or file
	Config = &ConfigTemplate{
		PKIRoot:    "./pki",
		LocalKey:   "client.key",
		ServerKey:  "server.key",
		AuthSigStr: "authenticate",
	}

	// setup crypto library
	rwcrypto.Init(Config.PKIRoot)
}
