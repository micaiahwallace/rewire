package rwcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// Keys is a package level keystore instance
var Keys Keystore

// Init package level crypto dependencies
func Init(pkipath string) {
	Keys = Keystore{PKIRoot: pkipath}
}

// GenerateKeyID hashes the public key as a uniuqe identifier
func GenerateKeyID(key *rsa.PublicKey) (string, error) {

	// Convert private key to byte slice
	bytes, err := GetKeyBytes(key, false)
	if err != nil {
		return "", err
	}

	// SHA256 hash the public key
	hashed := sha256.Sum256(bytes)

	// Convert byte slice to string
	return hex.EncodeToString(hashed[:]), nil
}

// GenerateKey returns a new generated private key
func GenerateKey(bits int) (*rsa.PrivateKey, error) {

	// Create keypair
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// CreateSignature returns msg signed with key as a string
func CreateSignature(key *rsa.PrivateKey, msg string) (string, error) {

	// get msg bytes as sha256 hash
	message := []byte(msg)
	hashed := sha256.Sum256(message)

	// sign message
	signedBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", errors.New("Unable to create signature")
	}

	// return bytes
	return hex.EncodeToString(signedBytes), nil
}

// ValidateSignature ensures the signature string matches the public key
func ValidateSignature(key *rsa.PublicKey, signatureString, msg string) bool {

	// get signature bytes
	signature, err := hex.DecodeString(signatureString)
	if err != nil {
		return false
	}

	// get msg string bytes as sha256 hash
	message := []byte(msg)
	hashed := sha256.Sum256(message)

	// verify signature
	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false
	}
	return true
}

// EncryptBytes encrypts bytes with a public key
func EncryptBytes(bytes []byte, key *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, bytes)
}

// DecryptBytes decrypts bytes with a private key
func DecryptBytes(bytes []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, bytes)
}
