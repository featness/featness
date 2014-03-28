package api

import (
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORKING"))
}
