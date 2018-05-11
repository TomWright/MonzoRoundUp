package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tomwright/monzoroundup/transport/http/webhook"
	"github.com/tomwright/monzoroundup/token"
	"github.com/tomwright/monzoroundup/user"
	"github.com/tomwright/monzoroundup/transport/http/auth"
)

var (
	bindAddress string
	router      *httprouter.Router
)

func Init(address string, tokenModel token.Model, userModel user.Model) {
	bindAddress = address
	router = httprouter.New()

	auth.Init(router)
	auth.InjectDependencies(tokenModel, userModel)
	
	webhook.Init(router)
}
