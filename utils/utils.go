package utils

import (
	"crypto/rand"
	"net"

	"golang.org/x/crypto/curve25519"
)

// Structure holding the command line arguments
type Arguments struct {
	Subnet    net.IPNet
	Output    string
	Endpoint  string
	Port      int
	All       bool
	Peers     int
	Keepalive int
}

// Creates a new X25519 key pair. Returns the privateKey, publickey and any errors
func NewX25519pair() ([]byte, []byte, error) {
	seed := make([]byte, 32)
	rand.Read(seed)

	privateKey, err := curve25519.X25519(seed, curve25519.Basepoint)
	if err != nil {
		return make([]byte, 32), make([]byte, 32), err
	}

	rand.Read(seed)

	publicKey, err := curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return make([]byte, 32), make([]byte, 32), err
	}

	return privateKey, publicKey, err
}

// Creates a new X25519 psk. Returns the key and any errors if they exist
func Newx25519Key() ([]byte, error) {
	seed := make([]byte, 32)
	rand.Read(seed)

	key, err := curve25519.X25519(seed, curve25519.Basepoint)
	if err != nil {
		return make([]byte, 32), err
	}

	return key, nil
}
