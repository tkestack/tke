package main

import (
	"fmt"
	"log"
)

func main() {
	if err := GenerateValueChart(); err != nil {
		log.Fatalf("generate chart value yaml fail, error: %s", err.Error())
	}
	fmt.Printf("generate chart value yaml success!")
}
