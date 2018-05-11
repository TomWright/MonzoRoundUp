package auth

import (
	"github.com/julienschmidt/httprouter"
)

const (
	LoginEndpoint    = "/auth/login"
	RegisterEndpoint = "/auth/register"
	CallbackEndpoint = "/auth/verify"
)

func Init(router *httprouter.Router) {
	router.POST(LoginEndpoint, loginHandler)
	router.POST(RegisterEndpoint, registerHandler)
	router.GET(CallbackEndpoint, callbackHandler)
}
