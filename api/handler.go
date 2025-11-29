package api

import "net/http"

// IHandler is a unified interface for all HTTP handlers
type IHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
