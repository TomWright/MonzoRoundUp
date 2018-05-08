package main

import (
	"github.com/tomwright/monzoroundup"
	"os"
)

func main() {
	oauthClientID := os.Args[1]
	oauthClientSecret := os.Args[2]

	monzoroundup.Init("localhost:8080", oauthClientID, oauthClientSecret)
	monzoroundup.Listen()
}
