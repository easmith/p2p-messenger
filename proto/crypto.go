package proto

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"io"
	"log"
	"os"
)

var randSeedFile = "seed.key"

// TODO:
func LoadKey() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {

	file, e := os.Open(randSeedFile)

	if e != nil {
		if os.IsNotExist(e) {
			// TODO: create file
			file = SaveKey()
		} else if os.IsPermission(e) {
			panic(e)
		}
	}

	publicKey, privateKey, e = ed25519.GenerateKey(file)
	if e != nil {
		log.Printf("Error: %v", e.Error())
	}
	return
}

func SaveKey() *os.File {
	file, e := os.Create(randSeedFile)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	seed := make([]byte, 32)
	_, e = rand.Reader.Read(seed)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	_, e = file.Write(seed)
	if e != nil {
		log.Fatalf("Error: %v", e.Error())
	}

	return file
}

//

func CreateKeyPair() (publicKey [32]byte, privateKey [32]byte) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	copy(publicKey[:], pub[:])
	copy(privateKey[:], priv[:])

	//priv[0] &= 248
	//priv[31] &= 127
	//priv[31] |= 64

	fmt.Printf("was : %v %v\n", publicKey, privateKey)

	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	fmt.Printf("be  : %v %v\n", publicKey, privateKey)

	//base := [32]byte{9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	fmt.Printf("will: %v %v\n", publicKey, privateKey)

	return
}

func GenerateKey(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var pub, priv [32]byte
	var err error

	_, err = io.ReadFull(rand, priv[:])
	if err != nil {
		return nil, nil, err
	}

	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	curve25519.ScalarBaseMult(&pub, &priv)

	return &priv, &pub, nil
}

//
//func Sign(msg []byte, privateKey [32]byte) {
//	ed25519.Sign(&privateKey, msg)
//}

func GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	var priv, pub, secret *[32]byte

	priv = privKey.(*[32]byte)
	pub = pubKey.(*[32]byte)
	secret = new([32]byte)

	curve25519.ScalarMult(secret, priv, pub)
	return secret[:], nil
}

func TestCrypto() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	log.Printf("PublicKey: %v", hex.EncodeToString(publicKey))
	log.Printf("PrivateKey: %s", hex.EncodeToString(privateKey))

	msg := "hello, world"

	log.Printf("Envelope to sign: %s", msg)

	sign := ed25519.Sign(privateKey, []byte(msg))

	log.Printf("Sign of Envelope: %s", hex.EncodeToString(sign))

	verify := ed25519.Verify(publicKey, []byte(msg), sign)

	log.Printf("Verification Status: %v", verify)

}
