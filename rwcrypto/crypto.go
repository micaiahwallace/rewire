package rwcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"micaiahwallace/rewire/rwutils"
	"os"
)

// CreateSignature returns a signed string of message using key
func CreateSignature(key *rsa.PrivateKey, msg string) (string, error) {

	// get hashed string
	message := []byte(msg)
	hashed := sha256.Sum256(message)

	// sign message
	signedBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", errors.New("Unable to create signature")
	}
	return hex.EncodeToString(signedBytes), nil
}

// ValidateSignature ensures the signature string matches the public key
func ValidateSignature(key *rsa.PublicKey, signatureString, msg string) bool {

	// get public key and hashed string
	signature, _ := hex.DecodeString(signatureString)
	message := []byte(msg)
	hashed := sha256.Sum256(message)

	// verify signature
	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Println("Signature verification error", err)
		return false
	}
	return true
}

// ParsePrivateKey coverts pem bytes into rsa private key
func ParsePrivateKey(pembytes []byte) (*rsa.PrivateKey, error) {

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

// ParsePrivateKeyFile runs ParsePrivateKey via filename instead of pem string
func ParsePrivateKeyFile(filename string) (*rsa.PrivateKey, error) {

	// verify file existence
	if exists := rwutils.FileExists(filename); !exists {
		return nil, errors.New("File does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fileStr := string(file)

	// Parse keyfile
	return ParsePrivateKey([]byte(fileStr))

}

// ParsePublicKey coverts pem bytes into rsa public key
func ParsePublicKey(pembytes []byte) (*rsa.PublicKey, error) {

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

// ParsePublicKeyFile runs ParsePublicKey via filename instead of pem string
func ParsePublicKeyFile(filename string) (*rsa.PublicKey, error) {

	// verify file existence
	if exists := rwutils.FileExists(filename); !exists {
		return nil, errors.New("File does not exist")
	}

	// Read contents of file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fileStr := string(file)

	// Parse keyfile
	return ParsePublicKey([]byte(fileStr))

}

// CreateKeypair returns a new generated private key
func CreateKeypair(bits int) (*rsa.PrivateKey, error) {

	// default key size of 2048 bits
	if bits <= 0 {
		bits = 2048
	}

	// Create keypair
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// LoadKeypair loads rsa private key from file system
func LoadKeypair(filePath string) (*rsa.PrivateKey, error) {

	var key *rsa.PrivateKey
	privExists := rwutils.FileExists(filePath)

	// Check fs for server key, create if not exist
	if !privExists {

		// Create new keypair
		key, err := CreateKeypair(2048)
		if err != nil {
			return nil, err
		}

		// Export key to file
		err = ExportKeyFile(key, filePath, true)
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	// Load key from file
	key, err := ParsePrivateKeyFile(filePath)
	if err != nil {
		return nil, errors.New("Can't load key")
	}

	return key, nil
}

// ParsePEM returns a pem object for the given public / private key
func ParsePEM(key interface{}, private bool) (*pem.Block, error) {

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

// ExportKeyBytes gets the pem encoded public or private key as a string
func ExportKeyBytes(key interface{}, private bool) ([]byte, error) {

	// Parse pem block
	pblock, err := ParsePEM(key, private)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(pblock), nil
}

// ExportKeyFile saves a pem encoded key to a file, private or public specified by private bool
func ExportKeyFile(key interface{}, filename string, private bool) error {

	// Create or open file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	// Parse pem block
	pblock, err := ParsePEM(key, private)
	if err != nil {
		return err
	}

	// Save pem to file and close file
	err = pem.Encode(file, pblock)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// GenerateKeyID hashes the public key as a uniuqe identifier
func GenerateKeyID(key *rsa.PublicKey) (string, error) {
	bytes, err := ExportKeyBytes(key, false)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(bytes)
	return hex.EncodeToString(hashed[:]), nil
}
