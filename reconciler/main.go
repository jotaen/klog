package main

import (
	"fmt"
	"log"
	"main/reconciler"
)

func main() {
	entry, err := reconciler.Parse("")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", entry)
}
