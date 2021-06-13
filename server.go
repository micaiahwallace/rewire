package rewire

import (
	"micaiahwallace/rewire/rwcrypto"
)

// SetupServerKeys loads or generates required server key
func SetupServerKeys() error {

	// Check if local key exists
	if !rwcrypto.Keys.Contains(Config.LocalKey) {

		// generate new local key
		newKey, err := rwcrypto.GenerateKey(Config.KeyBitLength)
		if err != nil {
			return err
		}

		// save to the keystore
		rwcrypto.Keys.SavePrivate(Config.LocalKey, newKey)
	}

	return nil
}
