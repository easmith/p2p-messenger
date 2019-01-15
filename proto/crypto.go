package proto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"io"
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

	log.Printf("Key exchange pair %s %s", hex.EncodeToString(publicKey[:]), hex.EncodeToString(privateKey[:]))

	return
}

// Calculate shared secret
func CalcSharedSecret(publicKey []byte, privateKey []byte) (secret [32]byte) {
	var pubKey [32]byte
	var privKey [32]byte
	copy(pubKey[:], publicKey[:])
	copy(privKey[:], privateKey[:])

	curve25519.ScalarMult(&secret, &privKey, &pubKey)
	log.Printf("publicKey %s", hex.EncodeToString(pubKey[:]))
	log.Printf("privateKey %s", hex.EncodeToString(privKey[:]))
	log.Printf("SharedKey %s", hex.EncodeToString(secret[:]))
	return
}

//Encrypt
func Encrypt(content []byte, key []byte) []byte {
	tip := len(content) % aes.BlockSize
	if tip != 0 {
		repeat := bytes.Repeat([]byte("\x00"), aes.BlockSize-(tip))
		content = append(content, repeat...)
	}

	log.Printf("length whant to bee 0 = %v", len(content)%aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the encrypted.
	encrypted := make([]byte, aes.BlockSize+len(content))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], content)

	return encrypted
}

func Decrypt(encrypted []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(encrypted)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encrypted, encrypted)

	encrypted = bytes.Trim(encrypted, string([]byte("\x00")))

	return encrypted
}
