package handler

import "net/http"

type Ping struct {
}

func (h Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
