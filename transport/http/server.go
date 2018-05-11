package http

import (
	"net/http"
	"log"
)

func Listen() {
	log.Printf("HTTP server listening on `%s`", bindAddress)
	err := http.ListenAndServe(bindAddress, router)
	if err != nil {
		log.Fatal("http server error: ", err)
	}
}
