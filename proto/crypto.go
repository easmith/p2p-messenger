package proto

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"log"
)

func TestCrypto() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	log.Printf("PublicKey: %v", hex.EncodeToString(publicKey))
	log.Printf("PrivateKey: %s", hex.EncodeToString(privateKey))

	msg := "hello, world"

	log.Printf("Message to sign: %s", msg)

	sign := ed25519.Sign(privateKey, []byte(msg))

	log.Printf("Sign of Message: %s", hex.EncodeToString(sign))

	verify := ed25519.Verify(publicKey, []byte(msg), sign)

	log.Printf("Verification Status: %v", verify)

}
