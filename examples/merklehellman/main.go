package main

import (
	"fmt"
	"log"

	"github.com/crodriguezvega/go-knapsackcrypto/pkg/merklehellman"
)

func main() {
	privKey, pubKey, err := merklehellman.GenerateKeys(10000)
	if err != nil {
		log.Fatalln(err)
	}

	p := []byte("Squeamish Ossifrage")
	c, err := pubKey.Encrypt(p)

	if err != nil {
		log.Fatalln(err)
	}

	m, err := privKey.Decrypt(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(m[:len(p)]))
}
