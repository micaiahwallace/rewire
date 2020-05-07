package rwcrypto

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"micaiahwallace/rewire/rwutils"
	"path"
)

// Keystore manages keys in a filesystem
type Keystore struct {

	// Root file system path to pki directory
	PKIRoot string

	// Cache file system keys
	KeyCache map[string]interface{}
}

// Contains checks for a key at path
func (store *Keystore) Contains(kpath string) bool {

	// verify file existence
	return rwutils.FileExists(path.Join(store.PKIRoot, kpath))
}

// GetPrivate returns a private key at path
func (store *Keystore) GetPrivate(kpath string) (*rsa.PrivateKey, error) {

	// verify file existence
	if exists := store.Contains(kpath); !exists {
		return nil, errors.New("File does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(kpath)
	if err != nil {
		return nil, err
	}

	// Parse keyfile
	return ParsePrivatePEM(file)
}

// GetPublic returns a public key at path
func (store *Keystore) GetPublic(kpath string) (*rsa.PublicKey, error) {

	// verify file existence
	if exists := store.Contains(kpath); !exists {
		return nil, errors.New("File does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(kpath)
	if err != nil {
		return nil, err
	}

	// Parse keyfile
	return ParsePublicPEM(file)
}

// SavePrivate saves a private key to a path
func (store *Keystore) SavePrivate(kpath string, key *rsa.PrivateKey) error {
	return ExportKeyToFile(key, path.Join(store.PKIRoot, kpath), true)
}

// SavePublic saves a public key to a path
func (store *Keystore) SavePublic(kpath string, key *rsa.PublicKey) error {
	return ExportKeyToFile(key, path.Join(store.PKIRoot, kpath), false)
}
