package utils

import (
	"crypto/tls"
	"net/http"
)

func FindAndDelete(data []string, delete string) []string {
	var respuesta []string
	for i := 0; i < len(data); i++ {
		if data[i] != delete {
			respuesta = append(respuesta, data[i])
		}
	}
	return respuesta
}

func GetClientHTTPS() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
