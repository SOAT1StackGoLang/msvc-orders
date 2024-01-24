package transport

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHTTPHandler(svc any, logger any) http.Handler {
	r := mux.NewRouter()

	return r
}
