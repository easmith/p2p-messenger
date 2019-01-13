package proto

import (
	"fmt"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
)

func TestCreateKeyPair(t *testing.T) {
	publicKey1, privateKey1 := CreateKeyPair()

	publicKey2, privateKey2 := CreateKeyPair()

	fmt.Printf("1 pub %v\n", publicKey1)
	fmt.Printf("1 priv %v\n", privateKey1)
	fmt.Printf("2 pub %v\n", publicKey2)
	fmt.Printf("2 priv %v\n", privateKey2)

	var sec1 [32]byte
	var sec2 [32]byte

	curve25519.ScalarMult(&sec1, &privateKey2, &publicKey1)
	curve25519.ScalarMult(&sec2, &privateKey1, &publicKey2)

	t.Logf("Sec1 %v", sec1)
	t.Logf("Sec2 %v", sec2)

	equal := reflect.DeepEqual(sec1, sec2)

	fmt.Printf("Secrets are equals: %v", equal)
}

func TestCreateKeyPair2(t *testing.T) {
	publicKey1, privateKey1, _ := ed25519.GenerateKey(nil)

	publicKey2, privateKey2, _ := ed25519.GenerateKey(nil)

	fmt.Printf("1 pub %v\n", publicKey1)
	fmt.Printf("1 priv %v\n", privateKey1)
	fmt.Printf("2 pub %v\n", publicKey2)
	fmt.Printf("2 priv %v\n", privateKey2)

	sec1, _ := GenerateSharedSecret(publicKey1, privateKey2)
	sec2, _ := GenerateSharedSecret(publicKey2, privateKey1)

	t.Logf("Sec1 %v", sec1)
	t.Logf("Sec2 %v", sec2)

	equal := reflect.DeepEqual(sec1, sec2)

	fmt.Printf("Secrets are equals: %v", equal)
}

func TestSign(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)

	//privateKey, publicKey, _ := GenerateKey(rand.Reader)
	var privateKey2 [32]byte
	copy(privateKey2[:], privateKey)
	privateKey2[0] &= 248
	privateKey2[31] &= 127
	privateKey2[31] |= 64

	var pub, priv [32]byte
	copy(pub[:], publicKey[:])
	copy(priv[:], privateKey2[:])

	curve25519.ScalarBaseMult(&pub, &priv)

	message := []byte("Message")
	copy(privateKey, priv[:])
	sign := ed25519.Sign(privateKey, message)

	t.Logf("Pub: %v Sign %v", publicKey, sign)

	copy(publicKey, pub[:])

	verify := ed25519.Verify(publicKey, message, sign)
	t.Logf("Verify: %v", verify)
}
