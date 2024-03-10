package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	// go get github.com/gorilla/mux
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin)).Methods("POST")
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", withJwtAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store)).Methods("GET")
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer)).Methods("POST")

	log.Println("Starting server on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
