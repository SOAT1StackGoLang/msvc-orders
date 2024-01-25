package transport

import (
	"log"
	"net/http"
)

func NewHTTPServer(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatalf("failed listening and serving: %s", err)
	}
}
