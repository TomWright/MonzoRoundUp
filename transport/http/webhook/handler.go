package webhook

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const ReceiveEventEndpoint = "/monzo/event"

func Init(router *httprouter.Router) {
	router.POST(ReceiveEventEndpoint, handleEvent)
}

func handleEvent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Accepted"))
}
