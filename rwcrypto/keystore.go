package rwcrypto

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"micaiahwallace/rewire/rwutils"
	"os"
	"path"
	"path/filepath"
)

// Keystore manages keys in a filesystem
type Keystore struct {

	// Root file system path to pki directory
	PKIRoot string

	// Cache file system keys
	KeyCache map[string]interface{}
}

// GetKeyPath concats the pkiroot to a keyfile path
func (store *Keystore) GetKeyPath(kpath string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return path.Join(dir, store.PKIRoot, kpath)
}

// Contains checks for a key at path
func (store *Keystore) Contains(kpath string) bool {
	return rwutils.FileExists(store.GetKeyPath(kpath))
}

// SavePrivate saves a private key to a path
func (store *Keystore) SavePrivate(kpath string, key *rsa.PrivateKey) error {
	return ExportKeyToFile(key, store.GetKeyPath(kpath), true)
}

// SavePublic saves a public key to a path
func (store *Keystore) SavePublic(kpath string, key *rsa.PublicKey) error {
	return ExportKeyToFile(key, store.GetKeyPath(kpath), false)
}

// GetPrivate returns a private key at path
func (store *Keystore) GetPrivate(kpath string) (*rsa.PrivateKey, error) {

	// verify file existence
	if exists := store.Contains(kpath); !exists {
		return nil, errors.New("file does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(store.GetKeyPath(kpath))
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
		return nil, errors.New("file does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(store.GetKeyPath(kpath))
	if err != nil {
		return nil, err
	}

	// Parse keyfile
	return ParsePublicPEM(file)
}
