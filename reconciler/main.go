package main

import (
	"fmt"
	"log"
	"main/reconciler"
)

func main() {
	err, day := reconciler.Parse("")
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  fmt.Printf("--- t:\n%v\n\n", day)
}
