package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

//Genero el par de claves Privada y Publica
func GeneratePrivatePublicKeys() (*rsa.PrivateKey, rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	//La clave publica forma parte del struct *rsa.PrivateKey
	publicKey := privateKey.PublicKey
	return privateKey, publicKey
}

//Ponemos la clave privada en formato PEM y la pasamos a bytes para almacenarla
func PrivateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	return privateKeyPem
}

//Ponemos la clave publica en formato PEM y la pasamos a bytes para almacenarla
func PublicKeyToPem(publicKey *rsa.PublicKey) []byte {
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey),
		},
	)
	return publicKeyPem
}

//Le pasamos la clave privada como un []byte para obtenerla como una clave funcional de tipo *rsa.PrivateKey
func PemToPrivateKey(privateBytes []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateBytes)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey
}

//Le pasamos la clave publica como un []byte para obtenerla como una clave funcional de tipo *rsa.PublicKey
func PemToPublicKey(publicBytes []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicBytes)
	publicKey, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return publicKey
}

func CifrarRSA(publicKey *rsa.PublicKey, plainText []byte) []byte {
	cipherBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
	if err != nil {
		panic(err)
	}
	return cipherBytes
}

func DescifrarRSA(privateKey *rsa.PrivateKey, cipherText []byte) []byte {
	dataBytes, err := privateKey.Decrypt(nil, cipherText, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}
	return dataBytes
}
