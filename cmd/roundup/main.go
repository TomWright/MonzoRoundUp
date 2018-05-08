package main

import (
	"github.com/tomwright/monzoroundup"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Expected 3 arguments: oauthClientID, oauthClientSecret, bindAddress (e.g. localhost:8080)")
		os.Exit(0)
	}
	oauthClientID := os.Args[1]
	oauthClientSecret := os.Args[2]
	bindAddress := os.Args[3]

	monzoroundup.Init(bindAddress, oauthClientID, oauthClientSecret)
	monzoroundup.Listen()
}
