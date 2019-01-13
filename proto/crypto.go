package proto

import (
	"crypto/rand"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"log"
	"os"
)

// create new *.key file with random content
func SaveKey(fileName string) *os.File {
	file, e := os.Create(fileName)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	seed := make([]byte, 64)
	_, e = rand.Reader.Read(seed)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	_, e = file.Write(seed)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	_, e = file.Seek(0, 0)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	return file
}

func LoadKey(name string) (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {

	fileName := name + ".key"
	file, e := os.Open(fileName)

	if e != nil {
		if os.IsNotExist(e) {
			file = SaveKey(fileName)
		} else if os.IsPermission(e) {
			panic(e)
		}
	}

	publicKey, privateKey, e = ed25519.GenerateKey(file)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}
	return
}

// create pair for ECDHE
func CreateKeyExchangePair() (publicKey [32]byte, privateKey [32]byte) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	copy(publicKey[:], pub[:])
	copy(privateKey[:], priv[:])

	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	return
}

// Calculate shared secret
func CalcSharedSecret(publicKey [32]byte, privateKey [32]byte) (secret [32]byte) {
	curve25519.ScalarMult(&secret, &privateKey, &publicKey)
	return
}
