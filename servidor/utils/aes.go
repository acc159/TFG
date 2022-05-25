package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
)

//Cifrado AES Modo OFB
func CifrarAES(key []byte, iv []byte, plainText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//Preparamos la longuitud del []byte de ciphertext para poder almacenar el texto cifrado y tambien el IV al principio de este
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	//Cifro
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)
	//Almaceno el IV en los primeros 16 bytes del textoCifrado
	copy(ciphertext, iv)
	return ciphertext
}

//Descifrado AES Modo OFB
func DescifrarAES(key []byte, ciphertext []byte) []byte {
	iv := ciphertext[:aes.BlockSize]
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return plaintext
}

func GenerateKeyIV(pass []byte) ([]byte, []byte) {
	data := sha512.Sum512(pass)
	iv := data[:aes.BlockSize]     //16 Byte
	Kaes := data[aes.BlockSize:48] // 32 Byte
	return Kaes, iv
}
