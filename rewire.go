package rewire

import "micaiahwallace/rewire/rwcrypto"

// Config store configuration for rewire
var Config *ConfigTemplate

// InitLib inits shared rewire logic
func InitLib() {

	// load config from cli or file
	Config = &ConfigTemplate{
		PKIRoot:      "./pki",
		LocalKey:     "client.key",
		ServerKey:    "server.key",
		AuthSigStr:   "authenticate",
		KeyBitLength: 2048,
	}

	// setup crypto library
	rwcrypto.Init(Config.PKIRoot)
}
