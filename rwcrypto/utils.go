package rwcrypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// GetKeyBlockPEM returns a pem object for the given public / private key
func GetKeyBlockPEM(key interface{}, private bool) (*pem.Block, error) {

	// Create pem block
	var pemBlock *pem.Block

	if private {
		// Create from private key
		privKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("Cannot assert private key")
		}
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		}

	} else {
		// Create from public key
		pubKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("Cannot assert public key")
		}
		pemBlock = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubKey),
		}
	}

	return pemBlock, nil
}

// GetKeyBytes converts a key to a byte slice
func GetKeyBytes(key interface{}, private bool) ([]byte, error) {

	// Parse pem block
	pblock, err := GetKeyBlockPEM(key, private)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(pblock), nil
}

// ParsePrivatePEM converts pem bytes into rsa private key
func ParsePrivatePEM(pembytes []byte) (*rsa.PrivateKey, error) {

	// decode pem string
	privkeyBlock, _ := pem.Decode(pembytes)
	if privkeyBlock == nil {
		return nil, errors.New("Unable to decode private key pem")
	}

	// convert pem bytes into pkcs1 public key
	privkeyBytes := privkeyBlock.Bytes
	privkey, err := x509.ParsePKCS1PrivateKey(privkeyBytes)
	if err != nil {
		return nil, errors.New("Unable to parse PKCS1 private key")
	}
	return privkey, nil
}

// ParsePublicPEM coverts pem bytes into rsa public key
func ParsePublicPEM(pembytes []byte) (*rsa.PublicKey, error) {

	// decode pem string
	pubkeyBlock, _ := pem.Decode(pembytes)
	if pubkeyBlock == nil {
		return nil, errors.New("Unable to decode public key pem")
	}

	// convert pem bytes into pkcs1 public key
	pubkeyBytes := pubkeyBlock.Bytes
	pubkey, err := x509.ParsePKCS1PublicKey(pubkeyBytes)
	if err != nil {
		return nil, errors.New("Unable to parse PKCS1 public key")
	}
	return pubkey, nil
}

// ExportKeyToFile saves a pem encoded key to a file, private or public specified by private bool
func ExportKeyToFile(key interface{}, filename string, private bool) error {

	// Create or open file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse pem block
	pblock, err := GetKeyBlockPEM(key, private)
	if err != nil {
		return err
	}

	// Save pem to file and close file
	err = pem.Encode(file, pblock)
	if err != nil {
		return err
	}

	return nil
}
