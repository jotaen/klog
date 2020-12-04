package main

import (
	"fmt"
	"log"
	"main/parser"
)

func main() {
	entry, err := parser.Parse("")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", entry)
}
