package main

import (
	"fmt"
	"log"

	"tkestack.io/tke/hack/lightweight-install/installer"
)

func main() {
	if err := installer.GenerateValueChart(); err != nil {
		log.Fatalf("generate chart value yaml fail, error: %s", err.Error())
	}
	fmt.Printf("generate chart value yaml success!")
}
