package main

import (
	"fmt"
	"log"
	"net/http"
)

func initRoutes() {
	log.Println("[INFO] Initializing routes")
	http.Handle("/stream", methodRouter{
		"POST": http.HandlerFunc(addChannel),
		"GET":  http.HandlerFunc(listChannels),
	})
	http.Handle("/stream/", methodRouter{
		"GET":    http.HandlerFunc(serveChannel),
		"DELETE": http.HandlerFunc(delChannel),
	})
}

type methodRouter map[string]http.Handler

func (m methodRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if handler, ok := m[r.Method]; ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
