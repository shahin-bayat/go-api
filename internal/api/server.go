package api

import (
	"log"
	"net/http"
)

func (s *APIServer) Run() {
	// go get github.com/gorilla/mux
	// router := mux.NewRouter()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", makeHttpHandleFunc(s.handleLogin))
	mux.HandleFunc("GET /account", makeHttpHandleFunc(s.handleGetAccount))
	mux.HandleFunc("GET /account/{id}", withJwtAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store))
	mux.HandleFunc("POST /account", makeHttpHandleFunc(s.handleCreateAccount))
	mux.HandleFunc("DELETE /account/{id}", withJwtAuth(makeHttpHandleFunc(s.handleDeleteAccount), s.store))
	mux.HandleFunc("PUT /account/{id}", withJwtAuth(makeHttpHandleFunc(s.handleUpdateAccount), s.store))
	mux.HandleFunc("POST /transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Println("Starting server on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, mux))
}
