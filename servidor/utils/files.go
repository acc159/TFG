package utils

import (
	"fmt"
	"log"
	"os"
)

func ReadFile(fileName string) []byte {
	datos, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return datos
}

func writeFile(filename string, data []byte) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
