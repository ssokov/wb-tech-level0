package http

import "net/http"

func NewRouter(handler *OrderHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/order", handler.GetOrder)
	mux.Handle("/", http.FileServer(http.Dir("./web/")))
	return mux
}
