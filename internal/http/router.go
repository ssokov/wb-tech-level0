package http

import (
	"net/http"
)

func NewRouter(handler *OrderHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/order", handler.GetOrder)
	mux.HandleFunc("/ping/", handler.Ping)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/usr/local/share/static"))))

	return mux
}
