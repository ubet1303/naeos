package http

import "net/http"

type Router struct {
	handler Handler
}

func NewRouter() *Router {
	return &Router{handler: Handler{}}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
