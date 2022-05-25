package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

var AC *x509.Certificate

var ACprivateKey *rsa.PrivateKey

var ACpublicKey *rsa.PublicKey

var server_AC_Key string

func LoadAC() {
	server_AC_Key = os.Getenv("Server_AC_Key")
	if server_AC_Key == "" {
		fmt.Println("Inserta la clave para descifrar la clave privada de la AC")
		fmt.Scanf("%v\n", &server_AC_Key)
	}

	//Compruebo si existe el fichero del certificado de la AC y sino mando a crear el fichero junto a los de las claves publica y privada
	certFileAC := ReadFile("certs/ac/certificateAC.pem")
	if len(certFileAC) == 0 {
		CreateACcertificateAndKeys()
	}
	certFileAC = ReadFile("certs/ac/certificateAC.pem")
	privateKeyFileACcipher := ReadFile("certs/ac/privateKeyAC.pem")
	publicKeyFileAC := ReadFile("certs/ac/publicKeyAC.pem")
	AC = PemToCertificate(certFileAC)
	ACpublicKey = PemToPublicKey(publicKeyFileAC)
	//Clave privada primero necesito descifrarla

	key, _ := GenerateKeyIV([]byte(server_AC_Key))
	privateKeyFileAC := DescifrarAES(key, privateKeyFileACcipher)
	ACprivateKey = PemToPrivateKey(privateKeyFileAC)
}

func CreateACcertificateAndKeys() {
	//Creamos la autoridad certificadora
	ac := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"TFG"},
			Country:      []string{"ES"},
			Province:     []string{"Alicante"},
			Locality:     []string{"Elda"},
			PostalCode:   []string{"03600"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	//Creamos la clave publica y privada para el certificado
	acPrivateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	acPublicKey := &acPrivateKey.PublicKey

	//Generamos el certificado
	acBytes, err := x509.CreateCertificate(rand.Reader, ac, ac, &acPrivateKey.PublicKey, acPrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	os.MkdirAll("certs/ac/", os.ModePerm)
	//Escribir el certificado en un fichero

	acPEM := CertificateToPem(acBytes)
	writeFile("certs/ac/certificateAC.pem", acPEM)

	//Escribir la Clave publica en un fichero
	publicKeyPEM := PublicKeyToPem(acPublicKey)
	writeFile("certs/ac/publicKeyAC.pem", publicKeyPEM)

	//Con la clave privada primero la cifro con AES para guardarla de manera segura
	privateKeyPEM := PrivateKeyToPem(acPrivateKey)

	key, iv := GenerateKeyIV([]byte(server_AC_Key))

	privateKeyCipher := CifrarAES(key, iv, privateKeyPEM)

	//Escribir la Clave privada en un fichero
	writeFile("certs/ac/privateKeyAC.pem", privateKeyCipher)
}

func CreateUserCertificate(username string, userPublicKeyBytes []byte) {

	// Creamos una plantilla de certificado que en este caso no es el de una AC sino el de un usuario normal
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"TFG"},
			Country:      []string{"ES"},
			Province:     []string{"Alicante"},
			Locality:     []string{"Elda"},
			PostalCode:   []string{"03600"},
		},
		IsCA:         false,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	//Para crear el certificado para un usuario necesitamos la plantilla del certificado, el certificado de la AC, la clave publica del usuario y la clave privada de la AC
	userPublicKey := PemToPublicKey(userPublicKeyBytes)
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, AC, userPublicKey, ACprivateKey)
	if err != nil {
		fmt.Println(err)
	}

	os.MkdirAll("certs/users/", os.ModePerm)
	//Lo pongo en formato PEM y lo escribo en el fichero correspondiente
	certUserPEM := CertificateToPem(certBytes)
	writeFile("certs/users/"+username+"_cert.pem", certUserPEM)

}

func GetUserCertificate(username string) []byte {
	return ReadFile("certs/users/" + username + "_cert.pem")
}

func GetACpublicKey() []byte {
	return PublicKeyToPem(ACpublicKey)
}

//De PEM a certificado y claves

func PemToCertificate(certBytes []byte) *x509.Certificate {
	block, _ := pem.Decode(certBytes)
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	return cert
}

func PemToPublicKey(publicBytes []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicBytes)
	publicKey, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return publicKey
}

func PemToPrivateKey(privateBytes []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateBytes)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey
}

//De certificado y claves a PEM

func CertificateToPem(certificateBytes []byte) []byte {
	certificatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certificateBytes,
		},
	)
	return certificatePem
}

func PrivateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	return privateKeyPem
}

func PublicKeyToPem(publicKey *rsa.PublicKey) []byte {
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey),
		},
	)
	return publicKeyPem
}
