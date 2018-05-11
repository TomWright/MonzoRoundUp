package http

import (
	"net/http"
	"log"
)

func Listen() {
	err := http.ListenAndServe(bindAddress, router)
	if err != nil {
		log.Fatal("http server error: ", err)
	}
}
